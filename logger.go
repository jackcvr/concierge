package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	quiet   bool
	verbose bool
}

func (l Logger) PrintError(format string, args ...any) {
	if !l.quiet {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func (l Logger) Info(format string, args ...any) {
	if !l.quiet {
		slog.Info(format, args...)
	}
}

func (l Logger) Debug(format string, args ...any) {
	if !l.quiet {
		slog.Debug(format, args...)
	}
}

func (l Logger) Error(format string, args ...any) {
	if !l.quiet {
		slog.Error(format, args...)
	}
}

func (l Logger) InitSLogger(flags int) {
	if flags > 0 {
		log.SetFlags(flags)
	} else {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}
	level := slog.LevelInfo
	if l.verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
