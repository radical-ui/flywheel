package flutter

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/radical-ui/flywheel/dart_lib"
)

func (self *Flutter) Configure(dartLib *dart_lib.DartLib) error {
	if err := mapFile(path.Join(self.dir, "pubspec.yaml"), makePubspecUpdater(dartLib)); err != nil {
		return err
	}

	return nil
}

func makePubspecUpdater(dartLib *dart_lib.DartLib) func(string) string {
	return func(text string) string {
		newDependencies := strings.Join([]string{
			"dependencies:",
			pathDependency("controller", dartLib.ControllerPath()),
			pathDependency("objects", dartLib.ObjectsPath()),
			"  flutter:\n    sdk: flutter",
		}, "\n")

		currentDependenciesMatcher := regexp.MustCompile(`dependencies:[\S\s]*sdk: flutter`)
		return currentDependenciesMatcher.ReplaceAllString(text, newDependencies)
	}
}

func pathDependency(name string, path string) string {
	return fmt.Sprintf("  %s:\n    path: %s\n", name, path)
}

func mapFile(file string, fn func(string) string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	before := string(bytes)
	after := fn(before)

	if before == after {
		slog.Warn("mapped a file, but nothing changed", "file", file)
		return nil
	}

	return os.WriteFile(file, []byte(after), os.ModePerm)
}
