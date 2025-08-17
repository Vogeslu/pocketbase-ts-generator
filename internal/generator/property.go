package generator

import (
	"fmt"
	"strings"

	"github.com/Vogeslu/pocketbase-ts-generator/internal/cmd"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/pocketbase_api"
	"github.com/iancoleman/strcase"
)

type InterfacePropertyType int

const (
	IptString = iota
	IptNumber
	IptBoolean
	IptJson
	IptFile
	IptEnum
	IptRelation
	IptGeopoint
)

type InterfaceProperty struct {
	Name           string
	CollectionName string
	Optional       bool
	Type           InterfacePropertyType
	IsArray        bool
	Data           interface{}
}

type CollectionWithProperties struct {
	Collection *pocketbase_api.Collection
	Properties []*InterfaceProperty
}

type propertyFlags struct {
	relationAsString bool
	forceOptional    bool
}

func GetInterfacePropertyType(typeName string) InterfacePropertyType {
	switch typeName {
	case "number":
		return IptNumber
	case "bool":
		return IptBoolean
	case "select":
		return IptEnum
	case "json":
		return IptJson
	case "file":
		return IptFile
	case "relation":
		return IptRelation
	case "geoPoint":
		return IptGeopoint
	default:
		return IptString
	}
}

func (propertyType InterfacePropertyType) String() string {
	switch propertyType {
	case IptString:
		return "String"
	case IptNumber:
		return "Number"
	case IptBoolean:
		return "Boolean"
	case IptEnum:
		return "Enum"
	case IptJson:
		return "Json"
	case IptFile:
		return "File"
	case IptRelation:
		return "Relation"
	case IptGeopoint:
		return "GeoPoint"
	}

	return "Unknown"
}

func (property InterfaceProperty) String() string {
	var data = []string{
		property.Type.String(),
	}

	if property.Optional {
		data = append(data, "Optional")
	}

	if property.IsArray {
		data = append(data, "Array")
	}

	if property.Type == IptRelation {
		relationTo, ok := property.Data.(string)
		if !ok {
			relationTo = "unknown (object)"
		}

		data = append(data, fmt.Sprintf("Relation to %s", relationTo))
	}

	if property.Type == IptEnum {
		enumData := property.Data.([]string)

		data = append(data, fmt.Sprintf("Enum Data [%s]", strings.Join(enumData, ", ")))
	}

	return fmt.Sprintf("%s (%s)", property.Name, strings.Join(data, ", "))
}

func (property InterfaceProperty) GetTypescriptProperty(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	return fmt.Sprintf("%s: %s", property.getTypescriptName(generatorFlags, flags), property.getTypescriptTypeWithArray(flags))
}

func (property InterfaceProperty) getTypescriptType(flags propertyFlags) string {
	switch property.Type {
	case IptNumber:
		return "number"
	case IptBoolean:
		return "boolean"
	case IptJson:
		if property.Optional {
			return "object | null | \"\""
		} else {
			return "object"
		}
	case IptEnum:
		return strcase.ToCamel(fmt.Sprintf("%s_%s_%s", property.CollectionName, property.Name, "options"))
	case IptRelation:
		if flags.relationAsString {
			return "string"
		}

		relationTo, ok := property.Data.(string)
		if !ok {
			return "object"
		} else {
			return strcase.ToCamel(relationTo)
		}
	case IptGeopoint:
		return "{ lon: number, lat: number }"
	default:
		return "string"
	}
}

func (property InterfaceProperty) getTypescriptTypeWithArray(flags propertyFlags) string {
	tsType := property.getTypescriptType(flags)

	if property.IsArray {
		if property.Optional {
			return fmt.Sprintf("%s[]", tsType)
		} else {
			return fmt.Sprintf("[%s]", tsType)
		}
	}

	return tsType
}

func (property InterfaceProperty) getTypescriptName(generatorFlags *cmd.GeneratorFlags, flags propertyFlags) string {
	if property.Optional && generatorFlags.MakeNonRequiredOptional || flags.forceOptional {
		return fmt.Sprintf("%s?", property.Name)
	}

	return property.Name
}

func (collection CollectionWithProperties) GetTypescriptInterface(generatorFlags *cmd.GeneratorFlags) string {
	properties := make([]string, len(collection.Properties))
	var additionalTypes []string
	var expandedRelations []string

	for i, property := range collection.Properties {
		properties[i] = fmt.Sprintf("    %s;", property.GetTypescriptProperty(generatorFlags, propertyFlags{forceOptional: false, relationAsString: true}))

		if property.Type == IptEnum {
			additionalTypes = append(additionalTypes, property.getTypescriptEnum())
		}

		if property.Type == IptRelation {
			expandedRelations = append(expandedRelations, fmt.Sprintf("    %s;", property.GetTypescriptProperty(generatorFlags, propertyFlags{forceOptional: true, relationAsString: false})))
		}
	}

	if len(expandedRelations) > 0 {
		expandedRelations = append(expandedRelations, "    [key: string]: unknown;")

		expandedType := fmt.Sprintf("export interface %sExpanded {\n%s\n}", strcase.ToCamel(collection.Collection.Name), strings.Join(expandedRelations, "\n"))

		additionalTypes = append(additionalTypes, expandedType)

		expandedLine := fmt.Sprintf("    expand?: %sExpanded;", strcase.ToCamel(collection.Collection.Name))

		properties = append([]string{expandedLine}, properties...)
	} else {
		expandedLine := "    expand?: { [key: string]: unknown; };"

		properties = append([]string{expandedLine}, properties...)
	}

	prefix := strings.Join(additionalTypes, "\n\n")

	if prefix != "" {
		prefix += "\n\n"
	}

	return fmt.Sprintf("%sexport interface %s {\n%s\n}", prefix, strcase.ToCamel(collection.Collection.Name), strings.Join(properties, "\n"))
}

func (property InterfaceProperty) getTypescriptEnum() string {
	if property.Type != IptEnum {
		return ""
	}

	enumData := property.Data.([]string)
	enumName := strcase.ToCamel(fmt.Sprintf("%s_%s_%s", property.CollectionName, property.Name, "options"))

	enumList := make([]string, len(enumData))

	for i, enum := range enumData {
		enumList[i] = fmt.Sprintf("    %s = \"%s\"", strcase.ToCamel(enum), enum)
	}

	return fmt.Sprintf("export enum %s {\n%s\n}", enumName, strings.Join(enumList, ",\n"))
}
