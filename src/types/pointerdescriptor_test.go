package types

import "testing"

func TestSimplePointerOfStruct(t *testing.T) {
	packageName := "package"
	structName := "aStruct"
	aStruct := NewStructDescriptor(packageName, structName)
	var typeCast TypeDescriptor = aStruct

	pointerDescriptor := NewPointerDescriptor(&typeCast)

	simplePointerName := pointerDescriptor.GetSimpleName()
	if simplePointerName != (PointerDescriptorMarker + structName) {
		t.Errorf("simplePointerName is wrong, expacted '%s', actual '%s'", (PointerDescriptorMarker + structName), simplePointerName)
	}
}

func TestPointerPlaceholderReplacementWithMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	pointerDescriptor := NewPointerDescriptor(&typeCastPlaceholder)

	switch (*(pointerDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("pointer reftype isn't a placeholder: '%T'", (*(pointerDescriptor.RefType)))
	}

	structDescriptor := NewStructDescriptor(packageName, structName)
	var typeCastStruct TypeDescriptor = structDescriptor

	pointerDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(pointerDescriptor.RefType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced pointer reftype isn't a struct: '%T'", (*(pointerDescriptor.RefType)))
	}
}

func TestPointerPlaceholderReplacementWithNonMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	pointerDescriptor := NewPointerDescriptor(&typeCastPlaceholder)

	switch (*(pointerDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("pointer reftype isn't a placeholder: '%T'", (*(pointerDescriptor.RefType)))
	}

	altStructName := "notTheStruct"
	structDescriptor := NewStructDescriptor(packageName, altStructName)
	var typeCastStruct TypeDescriptor = structDescriptor

	pointerDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(pointerDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		// it shouldn't have replaced because the cannonical names dont match
		t.Errorf("replaced pointer reftype isn't a placeholder: '%T'", (*(pointerDescriptor.RefType)))
	}
}
