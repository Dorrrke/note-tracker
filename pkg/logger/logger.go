package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger //nolint:gochecknoglobals //singltone

var once sync.Once //nolint:gochecknoglobals //singltone

func Get(flags ...bool) zerolog.Logger {
	once.Do(func() {
		if flags[0] {
			logger = zerolog.New(os.Stdout).
				Level(zerolog.DebugLevel).
				With().
				Timestamp().
				Caller().
				Logger().
				Output(
					zerolog.ConsoleWriter{
						Out:        os.Stderr,
						TimeFormat: "2006-01-02 15:04:05",
						PartsOrder: []string{
							"level",
							"time",
							"caller",
							"message",
						},
					})
		} else {
			logger = zerolog.New(os.Stdout).
				Level(zerolog.InfoLevel).
				With().
				Timestamp().
				Caller().
				Logger()
		}
	})

	return logger
}
