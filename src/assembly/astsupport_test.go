package assembly

import (
	"go/ast"
	"testing"
)

func TestExtractFieldNameFromBlankField(t *testing.T) {
	field := ast.Field{}
	name := ExtractFieldName(&field)

	if name != "" {
		t.Errorf("name is wrong, expected '%s' was '%s'", "", name)
	}
}

func TestExtractFieldNameFromNilField(t *testing.T) {
	var field ast.Field
	name := ExtractFieldName(&field)

	if name != "" {
		t.Errorf("name is wrong, expected '%s' was '%s'", "", name)
	}
}

func TestExtractFieldNameFromList(t *testing.T) {
	field := ast.Field{}
	field.Names = make([]*ast.Ident, 0)

	ident := &ast.Ident{}
	ident.Name = "this is a test"

	field.Names = append(field.Names, ident)

	name := ExtractFieldName(&field)

	if name != ident.Name {
		t.Errorf("name is wrong, expected '%s' was '%s'", ident.Name, name)
	}
}
