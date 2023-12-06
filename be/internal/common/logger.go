package common

import "go.uber.org/zap"

var logger, _ = zap.NewProduction()

func L() *zap.SugaredLogger {
	return logger.Sugar()
}
