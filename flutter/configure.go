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
	if err := mapFile(path.Join(self.dir, "pubspec.yaml"), self.makePubspecUpdater()); err != nil {
		return err
	}

	mainDartGenerator, err := self.makeMainDartGenerator()
	if err != nil {
		return err
	}

	if err := mapFile(path.Join(self.dir, "lib", "main.dart"), mainDartGenerator); err != nil {
		return err
	}

	return nil
}

func (self *Flutter) makeMainDartGenerator() (func(string) string, error) {
	boilerplate := `
		import 'package:controller/controller.dart';
		import 'package:flutter/material.dart';
		%v
	
		void main() {
			runApp(const MyApp());
		}
	
		class MyApp extends StatelessWidget {
			const MyApp({super.key});

			@override
			Widget build(BuildContext context) {
				return Controller(url: '%v', builder: (anyObject) {
					var name = anyObject.getName();
					var attributes = anyObject.getAttributes();

					%s

					throw Exception('unknown object kind: $name');
				});
			}
		}
	`

	importString := "import 'package:objects/%s';"
	statementString := `if (name == "%s") { return %s(%s); }`

	var imports []string
	var statements []string

	dartObjects, err := self.dartDoc.GetObjects()
	if err != nil {
		return nil, err
	}

	for _, widget := range dartObjects {
		imports = append(imports, fmt.Sprintf(importString, widget.File))

		fields, err := self.dartDoc.GetClassFields(widget.Name)
		if err != nil {
			return nil, err
		}

		var attributes string
		for _, attribute := range fields {
			attributes = attributes + fmt.Sprintf("%s: attributes['%s'],", attribute.Name, attribute.Name)
		}

		statements = append(statements, fmt.Sprintf(statementString, widget.Name, widget.Name, attributes))
	}

	return func(_ string) string {
		return fmt.Sprintf(boilerplate, strings.Join(imports, "\n"), self.options.Url, strings.Join(statements, "\n"))
	}, nil
}

func (self *Flutter) makePubspecUpdater() func(string) string {
	return func(text string) string {
		newDependencies := strings.Join([]string{
			"dependencies:",
			pathDependency("controller", self.dartLib.ControllerPath()),
			pathDependency("objects", self.dartLib.ObjectsPath()),
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
