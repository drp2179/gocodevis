package types

import (
	"strings"
	"unicode"
)

// MethodDescriptor - defines a method
type MethodDescriptor struct {
	Name    string
	Params  []*AttributeDescriptor
	Returns []*AttributeDescriptor
}

// TypeWithMethods - a Type that has methods
type TypeWithMethods interface {
	TypeDescriptor
	AddMethod(method *MethodDescriptor)
	GetMethods() []*MethodDescriptor
}

//NewMethodDescriptor - create a new method descriptor
func NewMethodDescriptor(methodName string, params []*AttributeDescriptor, returns []*AttributeDescriptor) *MethodDescriptor {
	method := MethodDescriptor{}

	method.Name = methodName
	method.Params = params
	method.Returns = returns

	return &method
}

// GetSignature - returns the signature of the method
func (method *MethodDescriptor) GetSignature() string {
	sb := strings.Builder{}

	sb.WriteString("f(")

	for i, p := range method.Params {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(p.GetSignature())
	}

	sb.WriteString(")(")

	for i, r := range method.Returns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(r.GetSignature())
	}

	sb.WriteString("(")

	return sb.String()
}

// IsPublic - true if public, false if private
func (method MethodDescriptor) IsPublic() bool {
	firstCharByte := method.Name[0]
	firstCharRune := rune(firstCharByte)
	return !unicode.IsLower(firstCharRune)
}
