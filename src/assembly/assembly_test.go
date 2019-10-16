package assembly

import (
	"go/ast"
	"testing"
	"types"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStructGeneratesClass(t *testing.T) {
	master := NewMasterTypeCollection()
	packageName := "a-package-name"
	typeName := "testclassname"

	cannonicalName := packageName + "." + typeName

	typeSpec := ast.TypeSpec{}
	typeSpec.Name = &ast.Ident{}
	typeSpec.Name.Name = typeName
	typeSpec.Type = &ast.StructType{}

	master.UpdateFromTypeSpec(packageName, &typeSpec)

	class, exists := master.GetDataType(cannonicalName)

	assert.True(t, exists, "the cannonical mapping should exist")
	if t.Failed() {
		return
	}

	assert.NotNil(t, class, "the classes should not be nil")
	if t.Failed() {
		return
	}
	//assert.Equal(t, typeName, (*class).GetName(), "the name of the class is wrong")
}
func TestStructWithAttribute(t *testing.T) {
	master := NewMasterTypeCollection()
	packageName := "a-package-name"
	structTypeName := "testclassname"
	cannonicalName := packageName + "." + structTypeName

	structType := makeStructType()

	fieldName := "testfieldName"
	fieldTypeName := "string"
	addFieldToStructType(structType, fieldName, fieldTypeName)

	typeSpec := ast.TypeSpec{}
	typeSpec.Name = &ast.Ident{}
	typeSpec.Name.Name = structTypeName
	typeSpec.Type = structType

	master.UpdateFromTypeSpec(packageName, &typeSpec)

	classType, _ := master.GetDataType(cannonicalName)

	switch actualClass := (*classType).(type) {
	case types.StructDescriptor:
		assert.Equal(t, 1, len(actualClass.Attributes), "the num of class attributes is wrong")
		if t.Failed() {
			return
		}

		attribute := actualClass.Attributes[0]
		assert.Equal(t, fieldName, attribute.Name, "the name of the attribute is wrong")
		assert.Equal(t, "string", attribute.CannonicalTypeName, "the type name of the attribute is wrong")
	default:
		assert.Fail(t, "wrong type")
	}

}

// ----------------------------------------------------------------------
//
// SUPPORT FUNCTIONS
//
// ----------------------------------------------------------------------

func addFieldToStructType(structType *ast.StructType, fieldName string, fieldTypeName string) {
	field := ast.Field{}
	field.Names = make([]*ast.Ident, 0)

	fieldNameIdent := ast.Ident{}
	fieldNameIdent.Name = fieldName
	field.Names = append(field.Names, &fieldNameIdent)

	fieldType := makeTypeExprFor(fieldTypeName)
	field.Type = fieldType

	structType.Fields.List = append(structType.Fields.List, &field)
}

func makeTypeExprFor(typeName string) ast.Expr {

	switch typeName {
	case "string":
		expr := &ast.Ident{}
		expr.Name = "string"
		return expr
	}

	return nil
}

func makeStructType() *ast.StructType {
	structType := ast.StructType{}
	structType.Fields = &ast.FieldList{}
	return &structType
}
