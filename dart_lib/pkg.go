package dart_lib

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type DartLib struct {
	dir string
}

func NewDartLib() (*DartLib, error) {
	// if we have some objects in the current working directory, we will use those
	bytes, _ := os.ReadFile("objects/pubspec.yaml")
	if strings.HasPrefix(string(bytes), "name: objects\n") {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		return &DartLib{cwd}, nil
	}

	flywheelPath, ok := os.LookupEnv("FLYWHEEL_PATH")
	if !ok {
		return nil, fmt.Errorf("couldn't find dart libraries in the current directory, and no FLYWHEEL_PATH env var was set")
	}

	self := &DartLib{flywheelPath}

	return self, nil
}

func (self *DartLib) ObjectsPath() string {
	return path.Join(self.dir, "objects")
}

func (self *DartLib) ControllerPath() string {
	return path.Join(self.dir, "controller")
}
