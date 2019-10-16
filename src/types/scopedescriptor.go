package types

import (
	"strings"
)

// ScopeDescriptor - represents the scope of an assembly
type ScopeDescriptor struct {
	Name           string
	CannonicalName string
	Methods        []*MethodDescriptor
}

// NewScopeDescriptor - create and initialize a new ScopeDescripter
func NewScopeDescriptor(scopeName string) *ScopeDescriptor {
	scope := ScopeDescriptor{}
	scope.Name = scopeName
	scope.CannonicalName = scopeName

	scope.Methods = make([]*MethodDescriptor, 0)

	return &scope
}

func (scope ScopeDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(scope.GetCannonicalName())
	sb.WriteString(" is scope {\n")

	for _, method := range scope.Methods {
		if method.IsPublic() {
			sb.WriteString("\tpublic ")
		} else {
			sb.WriteString("\tprivate ")
		}

		sb.WriteString(method.Name)
		sb.WriteString(" (")
		for i, p := range method.Params {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(p.GetSignature())
		}
		sb.WriteString(")")
		sb.WriteString("(")
		for i, r := range method.Returns {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(r.GetSignature())
		}
		sb.WriteString(")\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}

//-------------------------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (scope ScopeDescriptor) GetSimpleName() string {
	return scope.Name
}

// GetCannonicalName - Type interface to get the cannonical name
func (scope ScopeDescriptor) GetCannonicalName() string {
	return scope.CannonicalName
}

//-------------------------------------------
// TypeWithMethods interface
//

// AddMethod - add the method to the scope
func (scope *ScopeDescriptor) AddMethod(method *MethodDescriptor) {
	scope.Methods = append(scope.Methods, method)
}

// GetMethods - get the methods of the scope
func (scope *ScopeDescriptor) GetMethods() []*MethodDescriptor {
	return scope.Methods
}
