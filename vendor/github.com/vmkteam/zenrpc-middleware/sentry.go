package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/vmkteam/appkit"
	"github.com/vmkteam/zenrpc/v2"
)

// sentryContext is the context the call's details travel under on a Sentry
// event. Not extras: sentry-go deprecated them in v0.44 and removed them in
// v0.46. Not attributes either, their replacement: those reach logs and metrics
// only, never an error event.
const sentryContext = "rpc"

// WithSentry sets additional parameters for current Sentry scope. Context "rpc": params, duration, ip. Tags: platform,
// version, method. It's also handles panic.
func WithSentry(serverName string) zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) (response zenrpc.Response) { //nolint:nonamedreturns // the deferred recover answers a panic through it
			defer func() {
				var err error
				var rec any
				if rec = recover(); rec != nil {
					switch e := rec.(type) {
					case error:
						err = e
					default:
						err = fmt.Errorf("%v", e)
					}

					response = zenrpc.NewResponseError(nil, zenrpc.InternalError, "", nil)
				}

				if hub := sentry.GetHubFromContext(ctx); hub != nil {
					start, platform, version, ip, xRequestID := time.Now(), appkit.PlatformFromContext(ctx), appkit.VersionFromContext(ctx), appkit.IPFromContext(ctx), appkit.XRequestIDFromContext(ctx)

					methodName := fullMethodName(serverName, zenrpc.NamespaceFromContext(ctx), method)

					hub.Scope().SetContext(sentryContext, sentry.Context{
						"params":   params,
						"duration": time.Since(start).String(),
						"ip":       ip,
					})

					hub.Scope().SetTags(map[string]string{
						"platform":   platform,
						"version":    version,
						"method":     methodName,
						"xRequestId": xRequestID,
					})

					if err != nil {
						hub.CaptureException(err)
					}
				}
			}()

			return h(ctx, method, params)
		}
	}
}

// WithErrorLogger logs all errors (ErrorCode==500 or < 0) via Printf func and sends them to Sentry. It also removes
// sensitive error data from response. It is good to use pkg/errors for stack trace support in sentry.
func WithErrorLogger(pf Printf, serverName string) zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			start, platform, version, ip, xRequestID := time.Now(), appkit.PlatformFromContext(ctx), appkit.VersionFromContext(ctx), appkit.IPFromContext(ctx), appkit.XRequestIDFromContext(ctx)
			namespace := zenrpc.NamespaceFromContext(ctx)

			r := h(ctx, method, params)
			if r.Error != nil && (r.Error.Code == http.StatusInternalServerError || r.Error.Code < 0) {
				duration := time.Since(start)
				methodName := fullMethodName(serverName, namespace, method)

				pf("ip=%s platform=%q version=%q method=%s duration=%v params=%s xRequestId=%q err=%q", ip, platform, version, methodName, duration, params, xRequestID, r.Error)

				captureRPCError(ctx, r.Error, sentry.Context{
					"params":     params,
					"duration":   duration.String(),
					"ip":         ip,
					"error.data": r.Error.Data,
					"error.code": r.Error.Code,
				}, map[string]string{
					"platform":   platform,
					"version":    version,
					"method":     methodName,
					"xRequestId": xRequestID,
				})

				// remove sensitive error data from response
				r.Error.Err = nil
				r.Error.Message = "Internal error"
			}

			return r
		}
	}
}

// WithErrorSLog logs all errors (ErrorCode==500 or < 0) via [slog.ErrorContext] func and sends them to Sentry. It also removes
// sensitive error data from response. It is good to use pkg/errors for stack trace support in sentry.
func WithErrorSLog(pf Print, serverName string, fn LogAttrs) zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			r := h(ctx, method, params)

			// get additional args, check for ErrSkipLog
			var args []any
			if fn != nil {
				args = fn(ctx, r)
				if len(args) == 1 {
					if e, ok := args[0].(error); ok && errors.Is(e, ErrSkipLog) {
						return r
					}
				}
			}

			if r.Error != nil && (r.Error.Code == http.StatusInternalServerError || r.Error.Code < 0) {
				start, platform, version, ip, xRequestID := time.Now(), appkit.PlatformFromContext(ctx), appkit.VersionFromContext(ctx), appkit.IPFromContext(ctx), appkit.XRequestIDFromContext(ctx)
				namespace := zenrpc.NamespaceFromContext(ctx)

				duration := time.Since(start)
				methodName := fullMethodName(serverName, namespace, method)

				t := time.Since(start)
				logArgs := append(additionalArgs(ctx), []any{
					"method", fullMethodName(serverName, zenrpc.NamespaceFromContext(ctx), method),
					"duration", t.String(),
					"durationMS", t.Milliseconds(),
					"params", params,
					"err", r.Error,
					"userAgent", appkit.UserAgentFromContext(ctx),
					"xRequestId", appkit.XRequestIDFromContext(ctx),
				}...)

				pf(ctx, "rpc error", append(logArgs, args...)...)

				captureRPCError(ctx, r.Error, sentry.Context{
					"params":     params,
					"duration":   duration.String(),
					"ip":         ip,
					"error.data": r.Error.Data,
					"error.code": r.Error.Code,
				}, map[string]string{
					"platform":   platform,
					"version":    version,
					"method":     methodName,
					"xRequestId": xRequestID,
				})

				// remove sensitive error data from response
				r.Error.Err = nil
				r.Error.Message = "Internal error"
			}

			return r
		}
	}
}

// captureRPCError reports an RPC error to Sentry with the call's details: to the
// request's hub when the context carries one, to the current hub otherwise. The
// details go into a scope of their own, the one the event is captured with. A
// fresh scope beside the current hub used to be filled and never captured with,
// so without a hub in the context the event went out bare.
func captureRPCError(ctx context.Context, rpcErr *zenrpc.Error, details sentry.Context, tags map[string]string) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetContext(sentryContext, details)
		scope.SetTags(tags)
		hub.CaptureException(rpcErr)
	})
}
