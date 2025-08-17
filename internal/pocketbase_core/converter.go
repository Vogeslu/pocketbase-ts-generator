package pocketbase_core

import (
	"github.com/Vogeslu/pocketbase-ts-generator/internal/pocketbase_api"
	"github.com/pocketbase/pocketbase/core"
)

func convertPBCollections(pbCollections []*core.Collection) []pocketbase_api.Collection {
	output := make([]pocketbase_api.Collection, len(pbCollections))

	for i, pbCollection := range pbCollections {
		output[i] = convertPBCollection(pbCollection)
	}

	return output
}

func convertPBCollection(pbCollection *core.Collection) pocketbase_api.Collection {
	return pocketbase_api.Collection{
		Id:     pbCollection.Id,
		Name:   pbCollection.Name,
		Type:   pbCollection.Type,
		System: pbCollection.System,
		Fields: convertPBFields(pbCollection.Fields),
	}
}

func convertPBFields(pbFields core.FieldsList) []pocketbase_api.CollectionField {
	var output []pocketbase_api.CollectionField

	for _, pbField := range pbFields {
		output = append(output, convertPBField(pbField))
	}

	return output
}

func convertPBField(pbField core.Field) pocketbase_api.CollectionField {
	field := pocketbase_api.CollectionField{
		Id:     pbField.GetId(),
		Name:   pbField.GetName(),
		Type:   pbField.Type(),
		Hidden: pbField.GetHidden(),
	}

	switch v := pbField.(type) {
	case *core.TextField:
		field.Required = v.Required
	case *core.EditorField:
		field.Required = v.Required
	case *core.NumberField:
		field.Required = v.Required
	case *core.BoolField:
		field.Required = v.Required
	case *core.EmailField:
		field.Required = v.Required
	case *core.URLField:
		field.Required = v.Required
	case *core.DateField:
		field.Required = v.Required
	case *core.GeoPointField:
		field.Required = v.Required
	case *core.SelectField:
		field.MaxSelect = v.MaxSelect
		field.Required = v.Required
		field.Values = v.Values
	case *core.FileField:
		field.MaxSelect = v.MaxSelect
		field.Required = v.Required
	case *core.RelationField:
		field.MaxSelect = v.MaxSelect
		field.Required = v.Required
		field.CollectionId = v.CollectionId
	}

	return field
}
