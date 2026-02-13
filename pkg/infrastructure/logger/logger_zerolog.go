package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Yolto7/api-candidates/pkg/domain/logger"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/utils"
	"github.com/rs/zerolog"
)

type ZeroLogLogger struct {
	logger zerolog.Logger
}

func NewZeroLogLogger() *ZeroLogLogger {
	zerolog.TimeFieldFormat = ""

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}

	zerolog.MessageFieldName = "message"

	zlogger := zerolog.New(os.Stdout).With().Logger()

	return &ZeroLogLogger{
		logger: zlogger,
	}
}

func (l *ZeroLogLogger) log(level logger.LogLevel, input any) {
	event := l.getEvent(level)

	switch v := input.(type) {
	case string:
		event.Msg(v)
	case map[string]any:
		msg, hasMsg := v["msg"]
		if hasMsg {
			delete(v, "msg")
			event.Fields(v).Msg(fmt.Sprintf("%v", msg))
		} else {
			event.Fields(v).Msg("")
		}
	case utils.SafeError:
		event.
			Str("message", v.Message).
			Str("error", fmt.Sprintf("%v", v.Error)).
			Str("stack", v.Stack).
			Msg(v.Message)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			event.Msgf("log: %v", v)
		} else {
			event.RawJSON("payload", b).Msg("")
		}
	}
}

func (l *ZeroLogLogger) getEvent(level logger.LogLevel) *zerolog.Event {
	switch level {
	case logger.DEBUG:
		return l.logger.Debug()
	case logger.INFO:
		return l.logger.Info()
	case logger.WARN:
		return l.logger.Warn()
	case logger.ERROR:
		return l.logger.Error()
	default:
		return l.logger.Info()
	}
}

func (l *ZeroLogLogger) Debug(input any) {
	l.log(logger.DEBUG, input)
}

func (l *ZeroLogLogger) Info(input any) {
	l.log(logger.INFO, input)
}

func (l *ZeroLogLogger) Warn(input any) {
	l.log(logger.WARN, input)
}

func (l *ZeroLogLogger) Error(input any) {
	l.log(logger.ERROR, input)
}
