package logger

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLogger struct {
	Logger zerolog.Logger
}

func getTimestamp() time.Time {
	return time.Now().UTC()
}

// func callerMarshal()

func NewZeroLogger() ZeroLogger {
	zerolog.TimestampFunc = getTimestamp
	zerolog.SetGlobalLevel(0)
	var l ZeroLogger
	l.Logger = log.Logger.With().Caller().Logger()
	return l
}

var Logger = NewZeroLogger()
