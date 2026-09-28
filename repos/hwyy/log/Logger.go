package log

import (
	"github.com/phsym/console-slog"
	"log/slog"
	"os"
)

var (
	Logger = slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
)
