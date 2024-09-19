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

func (l Logger) InitSLogger(file string, flags int) error {
	if flags > 0 {
		log.SetFlags(flags)
	} else {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}
	level := slog.LevelInfo
	if l.verbose {
		level = slog.LevelDebug
	}
	logger, err := NewSLogger(file, level)
	if err != nil {
		return err
	}
	slog.SetDefault(logger)
	return nil
}

func NewSLogger(file string, level slog.Level) (*slog.Logger, error) {
	w := os.Stdout
	if file != "" {
		var err error
		w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})), nil
}
