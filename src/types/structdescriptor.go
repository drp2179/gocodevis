package types

import (
	"strings"
)

const (
	// StructDescriptorMarker - marker for structs
	StructDescriptorMarker = "struct_"
)

// StructDescriptor - collection of data for a struct
type StructDescriptor struct {
	builtin       bool
	Package       string
	Name          string
	Attributes    []*AttributeDescriptor
	Methods       []*MethodDescriptor
	PromotedTypes []*TypeDescriptor
}

// NewStructDescriptor - create a new data type descriptor
func NewStructDescriptor(packageName string, dtName string) *StructDescriptor {
	dt := StructDescriptor{}
	dt.builtin = false
	dt.Package = packageName
	dt.Name = dtName

	dt.Attributes = make([]*AttributeDescriptor, 0)
	dt.Methods = make([]*MethodDescriptor, 0)
	dt.PromotedTypes = make([]*TypeDescriptor, 0)

	return &dt
}

// NewBuiltinStructDescriptor - new data type descriptor for a builtin type
func NewBuiltinStructDescriptor(dtName string) *StructDescriptor {
	dt := StructDescriptor{}
	dt.builtin = true
	dt.Package = ""
	dt.Name = dtName

	dt.Attributes = make([]*AttributeDescriptor, 0)
	dt.Methods = make([]*MethodDescriptor, 0)
	dt.PromotedTypes = make([]*TypeDescriptor, 0)

	return &dt
}

// String - string version of the class descripter
func (dt StructDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(dt.GetCannonicalName())

	if dt.IsBuiltin() {
		sb.WriteString(", a builtin type\n")
	} else {

		sb.WriteString(" is struct {\n")

		for _, promotedType := range dt.PromotedTypes {
			sb.WriteString("\t")
			sb.WriteString((*promotedType).GetCannonicalName())
			sb.WriteString("\n")
		}

		for _, attribute := range dt.Attributes {
			if attribute.IsPublic() {
				sb.WriteString("\tpublic ")
			} else {
				sb.WriteString("\tprivate ")
			}

			sb.WriteString(attribute.Name)
			sb.WriteString(" ")
			sb.WriteString(attribute.CannonicalTypeName)
			sb.WriteString("\n")
		}
		sb.WriteString("\t-----------------------------------\n")
		for _, method := range dt.Methods {
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
	}
	return sb.String()
}

// IsBuiltin - returns true if created as a builtin type
func (dt StructDescriptor) IsBuiltin() bool {
	return dt.builtin
}

// ------------------------------------------------
// Type interface
//

// GetSimpleName - get the name of the class
func (dt StructDescriptor) GetSimpleName() string {
	return dt.Name
}

// GetCannonicalName - returns package.name
func (dt StructDescriptor) GetCannonicalName() string {
	sb := strings.Builder{}

	if len(dt.Package) > 0 {
		sb.WriteString(dt.Package)
		sb.WriteString(".")
	}

	if len(dt.Name) > 0 {
		sb.WriteString(dt.Name)
	}

	if len(dt.Package) <= 0 && len(dt.Name) <= 0 {
		sb.WriteString("struct{")
		for i, a := range dt.Attributes {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(a.CannonicalTypeName)
		}
		sb.WriteString("}")
	}

	return sb.String()
}

// ---------------------------------------
// TypeWithAttributes
//

// AddAttribute - add atttribute to the class
func (dt *StructDescriptor) AddAttribute(attribute *AttributeDescriptor) {
	dt.Attributes = append(dt.Attributes, attribute)
}

// GetAttributes = get the attrbutes of the class
func (dt *StructDescriptor) GetAttributes() []*AttributeDescriptor {
	return dt.Attributes
}

// AddPromotedType - add a promoted type to the struct
func (dt *StructDescriptor) AddPromotedType(promotedType *TypeDescriptor) {
	dt.PromotedTypes = append(dt.PromotedTypes, promotedType)
}

// GetPromotedTypes - returns the promoted type collection
func (dt *StructDescriptor) GetPromotedTypes() []*TypeDescriptor {
	return dt.PromotedTypes
}

// ------------------------------------------------
// TypeWithMehtods interface
//

// AddMethod - add method to the class
func (dt *StructDescriptor) AddMethod(method *MethodDescriptor) {
	dt.Methods = append(dt.Methods, method)
}

// GetMethods - the collection of methods for this interface
func (dt *StructDescriptor) GetMethods() []*MethodDescriptor {
	return dt.Methods
}
