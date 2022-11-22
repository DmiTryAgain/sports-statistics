package message_handler

type Handler interface {
	Construct() Handler
	Destruct()
	HandleWithResponse(dto *Dto) (string, bool, error)
}
