package types

import "testing"

func TestSimpleArrayOfStruct(t *testing.T) {
	packageName := "package"
	structName := "aStruct"
	aStruct := NewStructDescriptor(packageName, structName)
	var typeCast TypeDescriptor = aStruct

	arrayDescriptor := NewArrayDescriptor(&typeCast)

	simpleArrayName := arrayDescriptor.GetSimpleName()
	if simpleArrayName != (ArrayDescriptorMarker + structName) {
		t.Errorf("simpleArrayName is wrong, expacted '%s', actual '%s'", (ArrayDescriptorMarker + structName), simpleArrayName)
	}
}

func TestArrayPlaceholderReplacementWithMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	arrayDescriptor := NewArrayDescriptor(&typeCastPlaceholder)

	switch (*(arrayDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("array reftype isn't a placeholder: '%T'", (*(arrayDescriptor.RefType)))
	}

	structDescriptor := NewStructDescriptor(packageName, structName)
	var typeCastStruct TypeDescriptor = structDescriptor

	arrayDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(arrayDescriptor.RefType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced array reftype isn't a struct: '%T'", (*(arrayDescriptor.RefType)))
	}
}

func TestArrayPlaceholderReplacementWithNonMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	arrayDescriptor := NewArrayDescriptor(&typeCastPlaceholder)

	switch (*(arrayDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("array reftype isn't a placeholder: '%T'", (*(arrayDescriptor.RefType)))
	}

	altStructName := "notTheStruct"
	structDescriptor := NewStructDescriptor(packageName, altStructName)
	var typeCastStruct TypeDescriptor = structDescriptor

	arrayDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(arrayDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		// it shouldn't have replaced because the cannonical names dont match
		t.Errorf("replaced array reftype isn't a placeholder: '%T'", (*(arrayDescriptor.RefType)))
	}
}
