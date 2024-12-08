package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/radical-ui/flywheel/dart_doc"
	"github.com/radical-ui/flywheel/dart_lib"
	"github.com/radical-ui/flywheel/flutter"
)

type runOptions struct {
	genBindings bool
	genFlutter  *flutter.Options
	preview     bool
}

func run(options runOptions) error {
	ctx := contextWithCliCancelation()

	dartLib, err := dart_lib.NewDartLib()
	if err != nil {
		return err
	}

	dartDoc, err := dart_doc.NewDartDoc(dartLib)
	if err != nil {
		return err
	}

	if options.genBindings {
		schema, err := dartDoc.GetSchema()
		if err != nil {
			return err
		}

		bindings, err := schema.GenBindings("flywheel")
		if err != nil {
			return err
		}

		os.Stdout.Write(bindings)
	}

	var flutterInstance *flutter.Flutter

	if options.genFlutter != nil {
		f, err := flutter.NewFlutter(dartLib, dartDoc, *options.genFlutter)
		if err != nil {
			return err
		}

		if err := f.Configure(dartLib); err != nil {
			return err
		}

		flutterInstance = f
	}

	if options.preview {
		if err := flutterInstance.Preview(ctx); err != nil {
			return err
		}
	}

	return nil
}

func runWithErrorHandling(logFile string, options runOptions) {
	if err := setupLogger(logFile); err != nil {
		fmt.Printf("failed to initiate logger: %s\n", err.Error())
	}

	if err := run(options); err != nil {
		fmt.Println(err)
	}
}

func contextWithCliCancelation() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 3)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan

		cancel()
	}()

	return ctx
}
