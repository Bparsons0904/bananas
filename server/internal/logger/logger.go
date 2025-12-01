package logger

import (
	"log/slog"
	"os"

	gologger "github.com/Bparsons0904/goLogger"
)

type Logger = gologger.Logger

func New(name string) Logger {
	config := gologger.Config{
		Name:      name,
		Format:    gologger.FormatText,
		Level:     slog.LevelInfo,
		Writer:    os.Stdout,
		AddSource: false,
	}
	return gologger.NewWithConfig(config)
}

func NewWithConfig(config gologger.Config) Logger {
	return gologger.NewWithConfig(config)
}
