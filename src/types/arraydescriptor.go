package types

import "strings"

const (
	// ArrayDescriptorMarker - array marker
	ArrayDescriptorMarker = "[]"
)

// ArrayDescriptor - represents an array of a class
type ArrayDescriptor struct {
	RefType *TypeDescriptor
}

// NewArrayDescriptor = create a new array descriptor
func NewArrayDescriptor(refType *TypeDescriptor) *ArrayDescriptor {
	array := ArrayDescriptor{}
	array.RefType = refType

	return &array
}

// GetArrayTypeNameForTypename - returns the array type name string for the given type name
func GetArrayTypeNameForTypename(typeName string) string {
	return ArrayDescriptorMarker + typeName
}

func (array ArrayDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(array.GetCannonicalName())
	sb.WriteString(" is array of ")

	sb.WriteString((*(array.RefType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// --------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (array ArrayDescriptor) GetSimpleName() string {
	return GetArrayTypeNameForTypename((*(array.RefType)).GetSimpleName())
}

// GetCannonicalName - Type interface to get the cannonical name
func (array ArrayDescriptor) GetCannonicalName() string {
	return GetArrayTypeNameForTypename((*(array.RefType)).GetCannonicalName())
}

// ReplaceRefTypeIfPlaceholder - replace the array's refernce type
func (array *ArrayDescriptor) ReplaceRefTypeIfPlaceholder(refType *TypeDescriptor) {
	switch (*(array.RefType)).(type) {
	case *PlaceholderDescriptor:
		arrayRefTypeName := (*(array.RefType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if arrayRefTypeName == newRefTypeName {
			array.RefType = refType
		}
	}
}
