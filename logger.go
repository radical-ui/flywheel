package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/radical-ui/flywheel/helpers"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

func setupLogger(logFile string) error {
	level := slog.LevelDebug

	if logFile == "" {
		log.SetOutput(io.Discard)

		return nil
	}

	if err := helpers.EnsureDirExists(path.Dir(logFile)); err != nil {
		return err
	}

	file, err := os.Create(logFile)
	if err != nil {
		return err
	}

	handler := tint.NewHandler(colorable.NewColorable(file), &tint.Options{
		AddSource: true,
		Level:     level,
	})

	slog.SetDefault(slog.New(handler))

	return nil
}
