package types

import (
	"fmt"
	"strings"
)

// PlaceholderDescriptor - a placeholder for a type to be resolved later
type PlaceholderDescriptor struct {
	Name    string
	Scope   string
	Methods []*MethodDescriptor
}

// NewPlaceholderDescriptor - create a new place holder for the given name and scope
func NewPlaceholderDescriptor(scope string, name string) *PlaceholderDescriptor {

	placeholder := PlaceholderDescriptor{}
	placeholder.Name = name
	placeholder.Scope = scope
	placeholder.Methods = make([]*MethodDescriptor, 0)

	return &placeholder
}

//-------------------------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (placeholder PlaceholderDescriptor) GetSimpleName() string {
	return placeholder.Name
}

// GetCannonicalName - Type interface to get the cannonical name
func (placeholder PlaceholderDescriptor) GetCannonicalName() string {
	sb := strings.Builder{}

	if len(placeholder.Scope) > 0 {
		sb.WriteString(placeholder.Scope)
		sb.WriteString(".")
	}

	if len(placeholder.Name) > 0 {
		sb.WriteString(placeholder.Name)
	}

	return sb.String()
}

func (placeholder PlaceholderDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(placeholder.GetCannonicalName())
	sb.WriteString(" is placeholder of struct/interface {")
	sb.WriteString("\n")

	if len(placeholder.Methods) > 0 {
		for _, m := range placeholder.Methods {
			if m.IsPublic() {
				sb.WriteString("\tpublic ")
			} else {
				sb.WriteString("\tprivate ")
			}
			sb.WriteString(fmt.Sprintf("%s\n", m.GetSignature()))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

//-------------------------------------------
// TypeWithMethods interface
//

// AddMethod - add a method to the placeholder
func (placeholder *PlaceholderDescriptor) AddMethod(method *MethodDescriptor) {
	placeholder.Methods = append(placeholder.Methods, method)
}

// GetMethods - get the methods from  the placeholder
func (placeholder PlaceholderDescriptor) GetMethods() []*MethodDescriptor {
	return placeholder.Methods
}
