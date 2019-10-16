package types

import (
	"fmt"
	"testing"
)

func TestSimpleNameIsNotLocalToEmptyScope(t *testing.T) {

	name := "simple"
	scope := ""
	resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(name, scope)

	if isLocal {
		t.Error("isLocal should be false")
	}

	if resolvedName != name {
		t.Error(fmt.Sprintf("resolvedName expected '%s' was '%s'", name, resolvedName))
	}

	if extractedPrefix != "" {
		t.Error(fmt.Sprintf("extractedPrefix expected '%s' was '%s'", "", extractedPrefix))
	}
}

func TestCompoundNameIsLocalToEmptyScope(t *testing.T) {

	name := "base.simple"
	scope := "base"
	resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(name, scope)

	if !isLocal {
		t.Error("isLocal should be true")
	}

	if resolvedName != "simple" {
		t.Error(fmt.Sprintf("resolvedName expected '%s' was '%s'", "simple", resolvedName))
	}

	if extractedPrefix != "" {
		t.Error(fmt.Sprintf("extractedPrefix expected '%s' was '%s'", "", extractedPrefix))
	}
}

func TestPointerNameIsLocalToEmptyScope(t *testing.T) {

	name := "base.simple"
	pointerName := GetPointerTypeNameForTypename(name)
	scope := "base"
	resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(pointerName, scope)

	if !isLocal {
		t.Error("isLocal should be true")
	}

	if resolvedName != "simple" {
		t.Error(fmt.Sprintf("resolvedName expected '%s' was '%s'", "simple", resolvedName))
	}

	if extractedPrefix != PointerDescriptorMarker {
		t.Error(fmt.Sprintf("extractedPrefix expected '%s' was '%s'", PointerDescriptorMarker, extractedPrefix))
	}
}
