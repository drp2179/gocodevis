package types

import "strings"

// AliasedTypeDescriptor - defines an aliased data type
type AliasedTypeDescriptor struct {
	Package     string
	Name        string
	AliasedType *TypeDescriptor
	Methods     []*MethodDescriptor
}

// NewAliasedTypeDescriptor - create a new aliased data type
func NewAliasedTypeDescriptor(packageName string, aliasName string, dataType *TypeDescriptor) *AliasedTypeDescriptor {
	aliasedDt := AliasedTypeDescriptor{}

	aliasedDt.Package = packageName
	aliasedDt.Name = aliasName
	aliasedDt.AliasedType = dataType

	aliasedDt.Methods = make([]*MethodDescriptor, 0)

	return &aliasedDt
}

// GetSimpleName - get the name of the alias
func (alias AliasedTypeDescriptor) GetSimpleName() string {
	return alias.Name
}

// GetCannonicalName - returns package.name
func (alias AliasedTypeDescriptor) GetCannonicalName() string {
	sb := strings.Builder{}

	if len(alias.Package) > 0 {
		sb.WriteString(alias.Package)
		sb.WriteString(".")
	}

	if len(alias.Name) > 0 {
		sb.WriteString(alias.Name)
	}

	return sb.String()
}

func (alias AliasedTypeDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(alias.GetCannonicalName())
	sb.WriteString(" is alias for ")

	sb.WriteString((*(alias.AliasedType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// ------------------------------------------------
// TypeWithMehtods interface
//

// AddMethod - add method to the class
func (alias *AliasedTypeDescriptor) AddMethod(method *MethodDescriptor) {
	alias.Methods = append(alias.Methods, method)
}

// GetMethods - the collection of methods for this interface
func (alias *AliasedTypeDescriptor) GetMethods() []*MethodDescriptor {
	return alias.Methods
}
