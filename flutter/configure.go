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
	mapper, err := self.makePubspecUpdater()
	if err != nil {
		return err
	}
	mapFile(path.Join(self.dir, "pubspec.yaml"), mapper)

	mapper, err = self.makeController()
	if err != nil {
		return err
	}
	mapFile(path.Join(self.dir, "main.dart"), mapper)

	return nil
}

func (self *Flutter) makeController() (func(string) string, error) {
	boilerplate := `
	import 'package:controller.dart';
	import 'package:objects/objects.dart';
	import 'package:flutter/material.dart';
	%v
	
	void main() {
		runApp(const MyApp())
	}
	
	class MyApp extents StatelessWidget {
		const MyApp({super.key})

		@override
		Widget build(BuildContext context) {
			return Controller(url: %v, builder: (json) {
				%v
			})
		}
	}`

	importString := "import package:%v.dart;"
	statementString := `if (kind == "%v") return %v(%v)`
	var imports []string
	var statements []string
	widgets, err := self.dartDoc.GetObjects()
	if err != nil {
		return nil, err
	}

	for _, widget := range widgets {
		var attributes string
		imports = append(imports, fmt.Sprintf(importString, widget.File))
		fields, err := self.dartDoc.GetClassFields(widget.Name)

		if err != nil {
			return nil, err
		}
		for _, attribute := range fields {
			attributes = attributes + fmt.Sprintf("%v: json['attributes']['%v'],", attribute.Name, attribute.Name)
		}
		statements = append(statements, fmt.Sprintf(statementString, widget.Name, widget.Name, attributes))
	}
	return func(_ string) string { return fmt.Sprintf(boilerplate, self.options.Url, statements) }, nil
}

func (self *Flutter) makePubspecUpdater() (func(string) string, error) {
	return func(text string) string {
		newDependencies := strings.Join([]string{
			"dependencies:",
			pathDependency("controller", self.dartLib.ControllerPath()),
			pathDependency("objects", self.dartLib.ObjectsPath()),
			"  flutter:\n    sdk: flutter",
		}, "\n")

		currentDependenciesMatcher := regexp.MustCompile(`dependencies:[\S\s]*sdk: flutter`)
		return currentDependenciesMatcher.ReplaceAllString(text, newDependencies)
	}, nil
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
