package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
)

//go:embed objects
var objectsFs embed.FS

//go:embed controller
var controllerFs embed.FS

func getDartLib() (string, error) {
	bytes, _ := os.ReadFile("objects/pubspec.yaml")
	if strings.HasPrefix(string(bytes), "name: objects\n") {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return cwd, nil
	}

	return writeDartLib()
}

func writeDartLib() (string, error) {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("expected to find a $HOME env var")
	}

	dartPath := path.Join(home, ".cache", "flywheel", "dart_lib")

	slog.Info("deterimined dart path; deleting and rewriting", "dartPath", dartPath)
	os.RemoveAll(dartPath)

	if err := os.MkdirAll(path.Join(dartPath, "objects"), os.ModePerm); err != nil {
		return dartPath, err
	}

	if err := writeEmbedFs(dartPath, "objects", objectsFs); err != nil {
		return dartPath, err
	}

	if err := os.MkdirAll(path.Join(dartPath, "controller"), os.ModePerm); err != nil {
		return dartPath, err
	}

	if err := writeEmbedFs(dartPath, "controller", controllerFs); err != nil {
		return dartPath, err
	}

	return dartPath, nil
}

func writeEmbedFs(root string, local string, fs embed.FS) error {
	entries, err := fs.ReadDir(local)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		newLocal := path.Join(local, entry.Name())

		if entry.IsDir() {
			if err := os.Mkdir(path.Join(root, newLocal), os.ModePerm); err != nil {
				return err
			}

			if err := writeEmbedFs(root, newLocal, fs); err != nil {
				return err
			}

			continue
		}

		osFile, err := os.Create(path.Join(root, newLocal))
		if err != nil {
			return errors.Join(fmt.Errorf("failed to open os file"), err)
		}

		embedFile, err := fs.Open(path.Join(newLocal))
		if err != nil {
			return errors.Join(fmt.Errorf("failed to open embed.Fs file"), err)
		}

		if _, err := io.Copy(osFile, embedFile); err != nil {
			return err
		}
	}

	return nil
}
