package types

import "strings"

// FunctionDescriptor - defines a function type
type FunctionDescriptor struct {
	Package string
	Name    string
	Params  []*AttributeDescriptor
	Returns []*AttributeDescriptor
}

// NewFunctionDescriptor - create a new function data type
func NewFunctionDescriptor(packageName string, functionName string) *FunctionDescriptor {
	functionDesc := FunctionDescriptor{}

	functionDesc.Package = packageName
	functionDesc.Name = functionName

	functionDesc.Params = make([]*AttributeDescriptor, 0)
	functionDesc.Returns = make([]*AttributeDescriptor, 0)

	return &functionDesc
}

func (function FunctionDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(function.GetCannonicalName())

	sb.WriteString(" is function of ")

	sb.WriteString("(")
	for i, p := range function.Params {
		if i > 0 {
			sb.WriteString(",")
		}
		if len(p.Name) > 0 {
			sb.WriteString(p.Name)
			sb.WriteString(" ")
		}
		sb.WriteString(p.CannonicalTypeName)
	}
	sb.WriteString(")")
	sb.WriteString("(")
	for i, r := range function.Returns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(r.CannonicalTypeName)
	}
	sb.WriteString(")")
	sb.WriteString("\n")

	return sb.String()
}

// GetSimpleName - get the name of the class
func (function FunctionDescriptor) GetSimpleName() string {
	sb := strings.Builder{}

	sb.WriteString("func ")
	sb.WriteString(function.Name)

	return sb.String()
}

// GetCannonicalName - returns package.name
func (function FunctionDescriptor) GetCannonicalName() string {
	sb := strings.Builder{}

	if len(function.Package) > 0 {
		sb.WriteString(function.Package)
		sb.WriteString(" ")
	}

	sb.WriteString("func")

	if len(function.Name) > 0 {
		if len(function.Package) > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(function.Name)
	}

	sb.WriteString("(")
	for i, p := range function.Params {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(p.CannonicalTypeName)
	}

	sb.WriteString(")(")

	for i, r := range function.Returns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(r.CannonicalTypeName)
	}

	sb.WriteString(")")

	return sb.String()
}

// // ------------------------------------------------
// // TypeWithAttributes interface
// //

// // AddAttribute - add attribute to the function
// func (alias *AliasedTypeDescriptor) AddAttribute(method *MethodDescriptor) {
// 	alias.Methods = append(alias.Methods, method)
// }

// // GetAttributes - the collection of methods for this interface
// func (alias *AliasedTypeDescriptor) GetAttributes() []*MethodDescriptor {
// 	return alias.Methods
// }

// AddParameter - add a parameter to the function definition
func (function *FunctionDescriptor) AddParameter(param *AttributeDescriptor) {
	function.Params = append(function.Params, param)
}

// AddReturn - add a return type to the function definition
func (function *FunctionDescriptor) AddReturn(retn *AttributeDescriptor) {
	function.Returns = append(function.Returns, retn)
}
