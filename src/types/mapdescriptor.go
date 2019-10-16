package types

import (
	"strings"
)

const (
	// MapDescriptorMarker - map descriptor
	MapDescriptorMarker = "map["
)

// MapDescriptor - represents an map type
type MapDescriptor struct {
	KeyType   *TypeDescriptor
	ValueType *TypeDescriptor
}

// NewMapDescriptor = create a new map descriptor
func NewMapDescriptor(keyType *TypeDescriptor, valueType *TypeDescriptor) *MapDescriptor {
	mapType := MapDescriptor{}
	mapType.KeyType = keyType
	mapType.ValueType = valueType

	return &mapType
}

// ReplaceKeyTypeIfPlaceholder - replace key
func (mapDesc *MapDescriptor) ReplaceKeyTypeIfPlaceholder(refType *TypeDescriptor) {
	switch (*(mapDesc.KeyType)).(type) {
	case *PlaceholderDescriptor:
		keyRefTypeName := (*(mapDesc.KeyType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if keyRefTypeName == newRefTypeName {
			mapDesc.KeyType = refType
		}
	}
}

// ReplaceValueTypeIfPlaceholder - replace value
func (mapDesc *MapDescriptor) ReplaceValueTypeIfPlaceholder(refType *TypeDescriptor) {
	switch (*(mapDesc.ValueType)).(type) {
	case *PlaceholderDescriptor:
		valueRefTypeName := (*(mapDesc.ValueType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if valueRefTypeName == newRefTypeName {
			mapDesc.ValueType = refType
		}
	}
}

// GetMapTypeNameForTypenames - returns the array type name string for the given type name
func GetMapTypeNameForTypenames(keyTypename string, valueTypeName string) string {
	return MapDescriptorMarker + keyTypename + "][" + valueTypeName + "]"
}

func (mapDesc MapDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString("Map of ")
	sb.WriteString((*(mapDesc.KeyType)).GetCannonicalName())
	sb.WriteString(" to ")
	sb.WriteString((*(mapDesc.ValueType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// --------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (mapDesc MapDescriptor) GetSimpleName() string {
	return GetMapTypeNameForTypenames((*(mapDesc.KeyType)).GetSimpleName(), (*(mapDesc.ValueType)).GetSimpleName())
}

// GetCannonicalName - Type interface to get the cannonical name
func (mapDesc MapDescriptor) GetCannonicalName() string {
	return GetMapTypeNameForTypenames((*(mapDesc.KeyType)).GetCannonicalName(), (*(mapDesc.ValueType)).GetCannonicalName())
}
