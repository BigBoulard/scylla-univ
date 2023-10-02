package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type logger struct {
	logger zerolog.Logger
	// see https://github.com/rs/zerolog#leveled-logging
}

type Logger interface {
	Debug(string, string, string)
	Info(string, string, string)
	Warn(string, string, string)
	Error(error, string, string, string)
	Fatal(error, string, string, string)
}

func NewLogger() Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zlog := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %s |", i)
			},
			FieldsExclude: []string{"host", "user_agent", "service"},
		}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
	return &logger{
		logger: zlog,
	}
}

func (l *logger) Debug(pkg string, method string, msg string) {
	l.logger.
		Debug().
		Str("method", method).
		Str("pkg", pkg).
		Msg(msg)
}

func (l *logger) Info(pkg string, method string, msg string) {
	l.logger.
		Info().
		Str("pkg", pkg).
		Str("method", method).
		Msg(msg)
}

func (l *logger) Warn(pkg string, method string, msg string) {
	l.logger.
		Warn().
		Str("pkg", pkg).
		Str("method", method).
		Msg(msg)
}

func (l *logger) Error(err error, pkg string, method string, msg string) {
	l.logger.Error().
		Err(err).
		Str("pkg", pkg).
		Str("method", method).
		Msg(msg)
}

func (l *logger) Fatal(err error, pkg string, method string, msg string) {
	l.logger.Fatal().
		Err(err).
		Str("pkg", pkg).
		Str("method", method).
		Msg(msg)
}
