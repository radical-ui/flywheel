package dart_doc

import (
	"log/slog"
	"path"
	"strings"
)

var FieldTypeString = 1
var FieldTypeNumber = 2
var FieldTypeBoolean = 3
var FieldTypeList = 4
var FieldTypeColor = 5
var FieldTypeBinding = 6
var FieldTypeOptional = 7
var FieldTypeRef = 10

type ClassField struct {
	Name        string
	Kind        int
	Description string
	TypeName    string
}

type DartObject struct {
	Name string
	File string
}

func (self *DartDoc) GetObjects() ([]DartObject, error) {
	var objects []DartObject

	for _, node := range self.nodes {
		if node.Kind != nodeKindClass {
			continue
		}

		fullPath := path.Join(self.apiDir, node.Href)
		docHtml, err := NewDocHtml(fullPath)
		if err != nil {
			return objects, err
		}

		slog.Debug("", "node", node, "base", docHtml.GetBaseClasses())

		for _, class := range docHtml.GetBaseClasses() {
			if class == "Widget" {
				objects = append(objects, DartObject{node.Name, node.EnclosedBy.Name + ".dart"})
				break
			}
		}
	}

	return objects, nil
}

func (self *DartDoc) GetClassFields(className string) ([]ClassField, error) {
	var fields []ClassField

	for _, node := range self.nodes {
		if node.Kind != nodeKindField {
			continue
		}

		if node.EnclosedBy.Name != className {
			continue
		}

		html, err := NewDocHtml(path.Join(self.apiDir, node.Href))
		if err != nil {
			return fields, err
		}

		typeHref, typeName, err := html.GetFieldTypeData()
		if err != nil {
			return fields, err
		}

		fields = append(fields, ClassField{Name: node.Name, Kind: getFieldKind(typeName, typeHref), Description: node.Desc})
	}

	return fields, nil
}

func getFieldKind(name string, href string) int {
	isExternal := strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://")

	if name == "String" && isExternal {
		return FieldTypeString
	}

	// TODO

	return FieldTypeRef
}
