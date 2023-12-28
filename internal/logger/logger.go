package logger

import "log/slog"

var Log *slog.Logger

func init() {
	// Initialize the Log variable
	Log = slog.Default()
}
