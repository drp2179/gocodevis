package types

import (
	"strings"
	"testing"
)

func TestMapSimpleArrayOfStruct(t *testing.T) {
	packageName := "package"
	keyStructName := "aStructKey"
	valueStructName := "aStructValue"

	keyStruct := NewStructDescriptor(packageName, keyStructName)
	var typeCastKey TypeDescriptor = keyStruct

	valueStruct := NewStructDescriptor(packageName, valueStructName)
	var typeCastValue TypeDescriptor = valueStruct

	mapDescriptor := NewMapDescriptor(&typeCastKey, &typeCastValue)

	simpleMapName := mapDescriptor.GetSimpleName()
	if !strings.HasPrefix(simpleMapName, MapDescriptorMarker) {
		t.Errorf("simpleMapName is wrong, should start with '%s', actual '%s'", MapDescriptorMarker, simpleMapName)
	}
}

func TestMapPlaceholderReplacementWithMatchingName(t *testing.T) {
	packageName := "package"
	keyStructName := "aStructKey"
	valueStructName := "aStructValue"

	keyPlaceholderDescriptor := NewPlaceholderDescriptor(packageName, keyStructName)
	valuePlaceholderDescriptor := NewPlaceholderDescriptor(packageName, valueStructName)

	var keyTypeCastPlaceholder TypeDescriptor = keyPlaceholderDescriptor
	var valueTypeCastPlaceholder TypeDescriptor = valuePlaceholderDescriptor

	mapDescriptor := NewMapDescriptor(&keyTypeCastPlaceholder, &valueTypeCastPlaceholder)

	switch (*(mapDescriptor.KeyType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("map key reftype isn't a placeholder: '%T'", (*(mapDescriptor.KeyType)))
	}

	switch (*(mapDescriptor.ValueType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("map value reftype isn't a placeholder: '%T'", (*(mapDescriptor.ValueType)))
	}

	keyStructDescriptor := NewStructDescriptor(packageName, keyStructName)
	valueStructDescriptor := NewStructDescriptor(packageName, valueStructName)
	var keyTypeCastStruct TypeDescriptor = keyStructDescriptor
	var valueTypeCastStruct TypeDescriptor = valueStructDescriptor

	mapDescriptor.ReplaceKeyTypeIfPlaceholder(&keyTypeCastStruct)

	switch (*(mapDescriptor.KeyType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced map key reftype isn't a struct: '%T'", (*(mapDescriptor.KeyType)))
	}

	mapDescriptor.ReplaceValueTypeIfPlaceholder(&valueTypeCastStruct)

	switch (*(mapDescriptor.ValueType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced map value reftype isn't a struct: '%T'", (*(mapDescriptor.ValueType)))
	}
}

func TestMapPlaceholderReplacementWithNonMatchingName(t *testing.T) {
	packageName := "package"
	keyStructName := "aStructKey"
	valueStructName := "aStructValue"

	keyPlaceholderDescriptor := NewPlaceholderDescriptor(packageName, keyStructName)
	valuePlaceholderDescriptor := NewPlaceholderDescriptor(packageName, valueStructName)

	var keyTypeCastPlaceholder TypeDescriptor = keyPlaceholderDescriptor
	var valueTypeCastPlaceholder TypeDescriptor = valuePlaceholderDescriptor

	mapDescriptor := NewMapDescriptor(&keyTypeCastPlaceholder, &valueTypeCastPlaceholder)

	switch (*(mapDescriptor.KeyType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("map key reftype isn't a placeholder: '%T'", (*(mapDescriptor.KeyType)))
	}

	switch (*(mapDescriptor.ValueType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("map value reftype isn't a placeholder: '%T'", (*(mapDescriptor.ValueType)))
	}

	altKeyStructName := "altkeystructname"
	altValueStructName := "altvaluestructname"
	altKeyStructDescriptor := NewStructDescriptor(packageName, altKeyStructName)
	altValueStructDescriptor := NewStructDescriptor(packageName, altValueStructName)
	var altKeyTypeCastStruct TypeDescriptor = altKeyStructDescriptor
	var altValueTypeCastStruct TypeDescriptor = altValueStructDescriptor

	mapDescriptor.ReplaceKeyTypeIfPlaceholder(&altKeyTypeCastStruct)

	switch (*(mapDescriptor.KeyType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("replaced map key reftype should be a placeholder: '%T'", (*(mapDescriptor.KeyType)))
	}

	mapDescriptor.ReplaceValueTypeIfPlaceholder(&altValueTypeCastStruct)

	switch (*(mapDescriptor.ValueType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("replaced map value reftype should be a placeholder: '%T'", (*(mapDescriptor.ValueType)))
	}
}
