package message_handler

import "sports-statistics/internal/service/dto"

type Dto struct {
	*dto.Dto
}

func (d *Dto) Construct(messageDto *dto.Dto) *Dto {
	d.Dto = messageDto

	return d
}
