package types

import "strings"

const (
	// PointerDescriptorMarker - marks a pointer
	PointerDescriptorMarker = "*"
)

// PointerDescriptor - represents a pointer to a Type
type PointerDescriptor struct {
	RefType *TypeDescriptor
}

// NewPointerDescriptor = create a new pointer descriptor
func NewPointerDescriptor(refType *TypeDescriptor) *PointerDescriptor {
	ptr := PointerDescriptor{}
	ptr.RefType = refType

	return &ptr
}

// GetPointerTypeNameForTypename - get the pointer type name for the pointer to given type name
func GetPointerTypeNameForTypename(typeName string) string {
	return PointerDescriptorMarker + typeName
}

func (ptr PointerDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(ptr.GetCannonicalName())
	sb.WriteString(" is pointer to ")

	sb.WriteString((*(ptr.RefType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// --------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (ptr PointerDescriptor) GetSimpleName() string {
	return GetPointerTypeNameForTypename((*(ptr.RefType)).GetSimpleName())
}

// GetCannonicalName - Type interface to get the cannonical name
func (ptr PointerDescriptor) GetCannonicalName() string {
	return GetPointerTypeNameForTypename((*(ptr.RefType)).GetCannonicalName())
}

// ReplaceRefTypeIfPlaceholder - replace the pointer type
func (ptr *PointerDescriptor) ReplaceRefTypeIfPlaceholder(refType *TypeDescriptor) {
	switch (*(ptr.RefType)).(type) {
	case *PlaceholderDescriptor:
		ptrRefTypeName := (*(ptr.RefType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if ptrRefTypeName == newRefTypeName {
			ptr.RefType = refType
		}
	}
}
