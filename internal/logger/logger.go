package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) {

	var logger *zap.Logger
	var err error

	if env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	Log = logger
}
