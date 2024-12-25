package interpreter

import (
	"pocketbase-ts-generator/internal/generator"
	"pocketbase-ts-generator/internal/pocketbase"
)

func InterpretCollections(collections []*pocketbase.Collection, allCollections []pocketbase.Collection) []*generator.CollectionWithProperties {
	output := make([]*generator.CollectionWithProperties, len(collections))

	for i, collection := range collections {
		output[i] = InterpretCollection(collection, allCollections)
	}

	return output
}

func InterpretCollection(collection *pocketbase.Collection, allCollections []pocketbase.Collection) *generator.CollectionWithProperties {
	output := &generator.CollectionWithProperties{
		Collection: collection,
	}

	for _, field := range collection.Fields {
		if field.Hidden {
			continue
		}

		output.Properties = append(output.Properties, InterpretProperty(field, collection, allCollections))
	}

	return output
}

func InterpretProperty(field pocketbase.CollectionField, collection *pocketbase.Collection, allCollections []pocketbase.Collection) *generator.InterfaceProperty {
	output := &generator.InterfaceProperty{
		Name:           field.Name,
		CollectionName: collection.Name,
		Type:           generator.GetInterfacePropertyType(field.Type),
		Optional:       !field.Required,
	}

	if output.Type == generator.IptEnum || output.Type == generator.IptRelation || output.Type == generator.IptFile {
		output.IsArray = field.MaxSelect > 1
	}

	if output.Type == generator.IptRelation {
		output.Data = nil

		for _, collection := range allCollections {
			if collection.Id == field.CollectionId {
				output.Data = collection.Name
				break
			}
		}
	}

	if output.Type == generator.IptEnum {
		data := make([]string, len(field.Values))

		for i, value := range field.Values {
			data[i] = value
		}

		output.Data = data
	}

	return output
}
