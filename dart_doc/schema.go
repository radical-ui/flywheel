package dart_doc

import (
	"fmt"
	"log/slog"

	objection_schema "github.com/radical-ui/objection/schema"
)

func (self *DartDoc) GetSchema() (*objection_schema.Schema, error) {
	schema := &objection_schema.Schema{}

	dartObjects, err := self.GetObjects()
	if err != nil {
		return nil, err
	}

	slog.Debug("object names", "names", dartObjects)

	for _, dartObject := range dartObjects {
		objectFields, err := self.GetClassFields(dartObject.Name)
		if err != nil {
			return schema, err
		}

		var structFields []objection_schema.ItemDef

		for _, field := range objectFields {
			schemaType, err := self.getSchemaType(field.Kind, field.TypeName)
			if err != nil {
				return schema, err
			}

			item := objection_schema.ItemDef{
				Name:        field.Name,
				Description: field.Description,
				SchemaType:  schemaType,
			}

			structFields = append(structFields, item)
		}

		object := objection_schema.ObjectDef{
			Name:       dartObject.Name,
			Attributes: objection_schema.SchemaType{Kind: "struct", Properties: structFields},
		}

		schema.Objects = append(schema.Objects, object)
	}

	return schema, nil
}

func (self *DartDoc) getSchemaType(fieldType int, _ string) (objection_schema.SchemaType, error) {
	if fieldType == FieldTypeString {
		return objection_schema.SchemaType{Kind: "string"}, nil
	}

	if fieldType == FieldTypeNumber {
		return objection_schema.SchemaType{Kind: "number"}, nil
	}

	if fieldType == FieldTypeBoolean {
		return objection_schema.SchemaType{Kind: "number"}, nil
	}

	// TODO the other types here
	// if fieldType == FieldTypeList {
	// 	return objection_schema.SchemaType{Kind: "list"}
	// }

	return objection_schema.SchemaType{}, fmt.Errorf("unknown field type '%d'", fieldType)
}
