package logger

import "log"

type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(loggerDto *Dto) {
	log.Printf("Message info: \n"+
		"UserId: %d\n"+
		"UserName: %s\n"+
		"ChatId: %d\n"+
		"ChatType: %s\n"+
		"Message: %s\n",
		loggerDto.GetUserId(),
		loggerDto.GetUserName(),
		loggerDto.GetChatId(),
		loggerDto.GetChatType(),
		loggerDto.GetText(),
	)
}
