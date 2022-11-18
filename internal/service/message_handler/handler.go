package message_handler

type Handler interface {
	Construct() Handler
	HandleWithResponse(dto *Dto) (string, bool, error)
}
