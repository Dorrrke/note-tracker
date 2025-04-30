package logger

import (
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger //nolint:gochecknoglobals //singltone

var once sync.Once //nolint:gochecknoglobals //singltone

func Get(flags ...bool) zerolog.Logger {
	once.Do(func() {
		zerolog.TimestampFieldName = "time"                                         //nolint:reassign //todo
		zerolog.LevelFieldName = "level"                                            //nolint:reassign //todo
		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string { //nolint:reassign //todo
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return file + ":" + strconv.Itoa(line)
		}
		if flags[0] {
			logger = zerolog.New(os.Stdout).
				Level(zerolog.DebugLevel).
				With().
				Timestamp().
				Caller().
				Logger().
				Output(zerolog.ConsoleWriter{Out: os.Stdout})
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
