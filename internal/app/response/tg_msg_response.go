package response

type Response struct {
	msg      string
	needSend bool
	err      error
}

func (r *Response) Construct(
	msg string,
	needSend bool,
	err error,
) *Response {
	r.msg = msg
	r.needSend = needSend
	r.err = err

	return r
}

func (r *Response) GetMessage() string {
	return r.msg
}

func (r *Response) GetNeedSend() bool {
	return r.needSend
}

func (r *Response) GetError() error {
	return r.err
}
