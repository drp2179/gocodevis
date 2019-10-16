package types

import "strings"

const (
	// EllipsisDescriptorMarker - marker of an ellipsis
	EllipsisDescriptorMarker = "..."
)

// EllipsisDescriptor - represents an ellipsis of a type
type EllipsisDescriptor struct {
	RefType *TypeDescriptor
}

// NewEllipsisDescriptor = create a new ellipsis descriptor
func NewEllipsisDescriptor(refType *TypeDescriptor) *EllipsisDescriptor {
	ellipsis := EllipsisDescriptor{}
	ellipsis.RefType = refType

	return &ellipsis
}

// GetEllipsisTypeNameForTypename - returns the array type name string for the given type name
func GetEllipsisTypeNameForTypename(typeName string) string {
	return EllipsisDescriptorMarker + typeName
}

func (ellipsis EllipsisDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(ellipsis.GetCannonicalName())
	sb.WriteString(" is ellipsis of ")

	sb.WriteString((*(ellipsis.RefType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// --------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (ellipsis EllipsisDescriptor) GetSimpleName() string {
	return GetEllipsisTypeNameForTypename((*(ellipsis.RefType)).GetSimpleName())
}

// GetCannonicalName - Type interface to get the cannonical name
func (ellipsis EllipsisDescriptor) GetCannonicalName() string {
	return GetEllipsisTypeNameForTypename((*(ellipsis.RefType)).GetCannonicalName())
}

// ReplaceRefTypeIfPlaceholder - replace the ellipsis' refernce type
func (ellipsis *EllipsisDescriptor) ReplaceRefTypeIfPlaceholder(refType *TypeDescriptor) {

	switch (*(ellipsis.RefType)).(type) {
	case *PlaceholderDescriptor:
		ellipsisRefTypeName := (*(ellipsis.RefType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if ellipsisRefTypeName == newRefTypeName {
			ellipsis.RefType = refType
		}
	}
}
