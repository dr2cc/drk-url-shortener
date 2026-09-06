package testlog

import (
	"log/slog"
	"strings"
	"testing"
)

// В тестах мы вынуждены объявлять заглушку логгера. Пусть пользу приносит!

type tWriter struct {
	t testing.TB
}

func (tw tWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}

// New создает *slog.Logger, который пишет в вывод тестов Go.
func New(t testing.TB) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(tWriter{t: t}, opts)
	return slog.New(handler)
}
