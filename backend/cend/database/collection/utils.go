package collection

import (
	"log/slog"
	"runtime"
)


func LogInfo(message string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		slog.Info("Could not retrieve caller information", "message", message)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	slog.Info(message, "file", file, "line", line, "function", funcName)
}