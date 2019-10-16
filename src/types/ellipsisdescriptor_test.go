package types

import "testing"

func TestEllipsisSimpleArrayOfStruct(t *testing.T) {
	packageName := "package"
	structName := "aStruct"
	aStruct := NewStructDescriptor(packageName, structName)
	var typeCast TypeDescriptor = aStruct

	ellipsisDescriptor := NewEllipsisDescriptor(&typeCast)

	simpleEllipsisName := ellipsisDescriptor.GetSimpleName()
	if simpleEllipsisName != (EllipsisDescriptorMarker + structName) {
		t.Errorf("simpleEllipsisName is wrong, expacted '%s', actual '%s'", (EllipsisDescriptorMarker + structName), simpleEllipsisName)
	}
}

func TestEllipsisPlaceholderReplacementWithMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	ellipsisDescriptor := NewEllipsisDescriptor(&typeCastPlaceholder)

	switch (*(ellipsisDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("ellipsis reftype isn't a placeholder: '%T'", (*(ellipsisDescriptor.RefType)))
	}

	structDescriptor := NewStructDescriptor(packageName, structName)
	var typeCastStruct TypeDescriptor = structDescriptor

	ellipsisDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(ellipsisDescriptor.RefType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced ellipsis reftype isn't a struct: '%T'", (*(ellipsisDescriptor.RefType)))
	}
}

func TestEllipsisPlaceholderReplacementWithNonMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	ellipsisDescriptor := NewEllipsisDescriptor(&typeCastPlaceholder)

	switch (*(ellipsisDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("ellipsis reftype isn't a placeholder: '%T'", (*(ellipsisDescriptor.RefType)))
	}

	altStructName := "notTheStruct"
	structDescriptor := NewStructDescriptor(packageName, altStructName)
	var typeCastStruct TypeDescriptor = structDescriptor

	ellipsisDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(ellipsisDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		// it shouldn't have replaced because the cannonical names dont match
		t.Errorf("replaced ellipsis reftype isn't a placeholder: '%T'", (*(ellipsisDescriptor.RefType)))
	}
}
