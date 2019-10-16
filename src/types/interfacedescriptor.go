package types

import (
	"strings"
)

// InterfaceDescriptor - the definition of an interface
type InterfaceDescriptor struct {
	Name          string
	Package       string
	Methods       []*MethodDescriptor
	PromotedTypes []*TypeDescriptor
}

// NewInterfaceDescriptor - create a new interface descripter
func NewInterfaceDescriptor(packageName string, name string) *InterfaceDescriptor {
	intfc := InterfaceDescriptor{}
	intfc.Package = packageName
	intfc.Name = name
	intfc.PromotedTypes = make([]*TypeDescriptor, 0)
	intfc.Methods = make([]*MethodDescriptor, 0)

	return &intfc
}

// AddPromotedType - add promoted type to interface
func (intfc *InterfaceDescriptor) AddPromotedType(promotedType *TypeDescriptor) {
	intfc.PromotedTypes = append(intfc.PromotedTypes, promotedType)
}

// String - returns the string-ized version of the interface
func (intfc InterfaceDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(intfc.GetCannonicalName())
	sb.WriteString(" is interface {\n")

	for _, promotedType := range intfc.PromotedTypes {
		sb.WriteString("\t")
		sb.WriteString((*promotedType).GetCannonicalName())
		sb.WriteString("\n")
	}

	for _, method := range intfc.Methods {
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
			sb.WriteString(p.Name)
			sb.WriteString(" ")
			sb.WriteString(p.CannonicalTypeName)
		}
		sb.WriteString(")")
		sb.WriteString("(")
		for i, r := range method.Returns {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(r.CannonicalTypeName)
		}
		sb.WriteString(")\n")
	}

	sb.WriteString("}\n")

	return sb.String()
}

// ------------------------------------------
// Type interface
//

// GetSimpleName - get the name of the class
func (intfc InterfaceDescriptor) GetSimpleName() string {
	return intfc.Name
}

// GetCannonicalName - returns package.name
func (intfc InterfaceDescriptor) GetCannonicalName() string {
	sb := strings.Builder{}

	if len(intfc.Package) > 0 {
		sb.WriteString(intfc.Package)
		sb.WriteString(".")
	}

	if len(intfc.Name) > 0 {
		sb.WriteString(intfc.Name)
	}

	if len(intfc.Package) <= 0 && len(intfc.Name) <= 0 {
		sb.WriteString("interface(")
		for i, m := range intfc.Methods {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(m.GetSignature())
		}
		sb.WriteString(")")
	}

	return sb.String()
}

// ------------------------------------------
// TypeWithMethods interface
//

// AddMethod - add method to the interface
func (intfc *InterfaceDescriptor) AddMethod(method *MethodDescriptor) {
	intfc.Methods = append(intfc.Methods, method)
}

// GetMethods - the collection of methods for this interface
func (intfc *InterfaceDescriptor) GetMethods() []*MethodDescriptor {
	return intfc.Methods
}
