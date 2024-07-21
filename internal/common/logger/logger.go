package logger

import (
	"go.uber.org/zap"
)

var ZapSugarLogger zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ZapSugarLogger = *logger.Sugar()
}
