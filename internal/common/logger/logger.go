package logger

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// CustomLogger is used to abstract away a 3rd party logger.
type CustomLogger struct{}

// Custom is used as a global var so I don't need to pass in everywhere as a dependency.
var Custom *CustomLogger

func (cl *CustomLogger) Debugln(msg ...interface{}) {
	log.Debug().Msg(makeMessage(msg))
}

func (cl *CustomLogger) Debugf(format string, msg ...interface{}) {
	log.Debug().Msg(makeMessage(fmt.Sprintf(format, msg...)))
}

func (cl *CustomLogger) Infoln(msg ...interface{}) {
	log.Info().Msg(makeMessage(msg))
}

func (cl *CustomLogger) Errorln(msg ...interface{}) {
	log.Error().Msg(makeMessage(msg))
}

func (cl *CustomLogger) Panicln(msg ...interface{}) {
	log.Panic().Msg(makeMessage(msg))
}

func makeMessage(msg ...interface{}) string {
	fullMessage := ""
	for _, m := range msg {
		fullMessage += fmt.Sprintf("%v ", m)
	}
	return fullMessage
}
