package types

import "testing"

func TestNewAliasedDataType(t *testing.T) {

	dtPackage := "dt-package"
	dtName := "dt-name"
	dt := NewStructDescriptor(dtPackage, dtName)

	var dtCast TypeDescriptor = dt

	aliasePackage := "alias-package"
	aliaseName := "aliase-name"
	adt := NewAliasedTypeDescriptor(aliasePackage, aliaseName, &dtCast)

	if adt == nil {
		t.FailNow()
	}

	if adt.Package != aliasePackage {
		t.FailNow()
	}
	if adt.Name != aliaseName {
		t.FailNow()
	}

	if adt.AliasedType != &dtCast {
		t.FailNow()
	}
}
