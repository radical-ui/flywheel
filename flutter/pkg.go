package flutter

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/radical-ui/flywheel/helpers"
)

type Options struct {
	BundleIdentifier string
	DisplayName      string
}

type Flutter struct {
	options Options
	dir     string
}

func NewFlutter(options Options) (*Flutter, error) {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		return nil, fmt.Errorf("failed to detect user home")
	}

	appsCacheDir := path.Join(home, ".cache", "flywheel", "apps")
	appCacheDir := path.Join(appsCacheDir, options.BundleIdentifier)

	if _, err := os.Stat(appCacheDir); err != nil {
		if os.IsNotExist(err) {
			if err := helpers.EnsureDirExists(appsCacheDir); err != nil {
				return nil, err
			}

			org, name, err := splitBundleIdentifier(options.BundleIdentifier)
			if err != nil {
				return nil, err
			}

			if err := createFlutterApp(appsCacheDir, options.BundleIdentifier, org, name); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &Flutter{options, appCacheDir}, nil
}

func createFlutterApp(appsCacheDir string, folderName string, org string, name string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command := exec.Command("flutter", "create", folderName, "--project-name", name, "--org", org)
	command.Dir = appsCacheDir
	command.Stderr = &stderr
	command.Stdout = &stdout

	err := errors.New("failed to run spinner for creating flutter project")

	slog.Info("running `flutter create`", "folderName", folderName, "cwd", appsCacheDir, "name", name, "org", org)

	spinner.New().
		Title("Creating flutter project").
		Action(func() {
			err = command.Run()
		}).
		Run()

	if err != nil {
		return errors.Join(fmt.Errorf("failed to create flutter project: %s", stderr.String()), err)
	}

	if stderr.Len() != 0 {
		slog.Warn("recieved stderr from `flutter create`", "stderr", stderr.String())
	}

	return nil
}

func splitBundleIdentifier(bundleIdentifier string) (string, string, error) {
	numberOfDots := strings.Count(bundleIdentifier, ".")
	if numberOfDots < 2 {
		return "", "", fmt.Errorf("the bundle identifier does not contain the correct number of dots")
	}

	lastDot := strings.LastIndex(bundleIdentifier, ".")

	return bundleIdentifier[0:lastDot], bundleIdentifier[lastDot+1:], nil
}
