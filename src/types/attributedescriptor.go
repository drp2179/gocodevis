package types

import "unicode"

// AttributeDescriptor - defines an attribute
type AttributeDescriptor struct {
	Name               string
	CannonicalTypeName string
	Tag                string
}

// TypeWithAttributes - a Type that has attributes
type TypeWithAttributes interface {
	TypeDescriptor

	AddAttribute(attribute *AttributeDescriptor)
	GetAttributes() []*AttributeDescriptor
	AddPromotedType(promotedType *TypeDescriptor)
	GetPromotedTypes() []*TypeDescriptor
}

// NewAttributeDescriptor - create a new attribute descriptor
func NewAttributeDescriptor(attributeName string, attributeType TypeDescriptor, tag string) *AttributeDescriptor {
	attribute := AttributeDescriptor{}

	attribute.Name = attributeName
	attribute.CannonicalTypeName = attributeType.GetCannonicalName()
	attribute.Tag = tag

	return &attribute
}

// GetSignature - get the signature of the attribute descriptor
func (attribute AttributeDescriptor) GetSignature() string {
	return attribute.CannonicalTypeName
}

// IsPublic - true if public, false if private
func (attribute AttributeDescriptor) IsPublic() bool {
	firstCharByte := attribute.Name[0]
	firstCharRune := rune(firstCharByte)
	return !unicode.IsLower(firstCharRune)
}
