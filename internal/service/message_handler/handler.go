package message_handler

type Handler interface {
	HandleWithResponse(dto *Dto) (string, bool)
}
