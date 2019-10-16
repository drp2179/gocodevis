package assembly

import (
	"go/ast"
	"strings"
	"testing"
)

func TestExtractIdentNoPriorDT(t *testing.T) {
	master := NewMasterTypeCollection()

	scope := "scope"
	ident := &ast.Ident{}
	ident.Name = "name"

	typeName, newType := master.ExtractIdentType(scope, ident)

	if !strings.HasPrefix(typeName, scope) {
		t.Errorf("typeName is wrong, scope is missing, expected starts with '%s' was '%s'", scope, typeName)
	}

	if !strings.HasSuffix(typeName, ident.Name) {
		t.Errorf("typeName is wrong, type is missing, expected ends with '%s' was '%s'", ident.Name, typeName)
	}

	if newType == nil {
		t.Error("newType should not be nil")
	}

	if _, placeholderExists := master.GetPlaceholder(typeName); !placeholderExists {
		t.Error("Should have created a placeholder")
	}
}

func TestExtractIdentPriorDT(t *testing.T) {
	master := NewMasterTypeCollection()

	scope := "scope"
	ident := &ast.Ident{}
	ident.Name = "name"

	typeName, newType := master.ExtractIdentType(scope, ident)

	if !strings.HasPrefix(typeName, scope) {
		t.Errorf("typeName is wrong, scope is missing, expected starts with '%s' was '%s'", scope, typeName)
	}

	if !strings.HasSuffix(typeName, ident.Name) {
		t.Errorf("typeName is wrong, type is missing, expected ends with '%s' was '%s'", ident.Name, typeName)
	}

	if newType == nil {
		t.Error("newType should not be nil")
	}

	if _, placeholderExists := master.GetPlaceholder(typeName); !placeholderExists {
		t.Error("Should have created a placeholder")
	}
}
