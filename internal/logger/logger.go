package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	// Initialize the Log variable
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil) // 👈
	Log = slog.New(jsonHandler)
}
