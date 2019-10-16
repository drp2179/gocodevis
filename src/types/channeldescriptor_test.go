package types

import (
	"strings"
	"testing"
)

func TestSimpleReceiveChannelOfStruct(t *testing.T) {
	packageName := "package"
	structName := "aStruct"
	aStruct := NewStructDescriptor(packageName, structName)
	var typeCast TypeDescriptor = aStruct

	channelDescriptor := NewChannelDescriptor(ChannelDirectionRx, &typeCast)

	simpleChannelName := channelDescriptor.GetSimpleName()
	if !strings.Contains(simpleChannelName, ChannelDescriptorMarker) {
		t.Errorf("simpleChannelName is wrong, does not contain channel marker '%s', actual '%s'", ChannelDescriptorMarker, simpleChannelName)
	}
	if !strings.Contains(simpleChannelName, ChannelDescriptorReceiveMarker) {
		t.Errorf("simpleChannelName is wrong, does not contain channel receive marker '%s', actual '%s'", ChannelDescriptorReceiveMarker, simpleChannelName)
	}
	if !strings.Contains(simpleChannelName, structName) {
		t.Errorf("simpleChannelName is wrong, does not contain structure  marker '%s', actual '%s'", structName, simpleChannelName)
	}
}

func TestSimpleTransmitChannelOfStruct(t *testing.T) {
	packageName := "package"
	structName := "aStruct"
	aStruct := NewStructDescriptor(packageName, structName)
	var typeCast TypeDescriptor = aStruct

	channelDescriptor := NewChannelDescriptor(ChannelDirectionTx, &typeCast)

	simpleChannelName := channelDescriptor.GetSimpleName()
	if !strings.Contains(simpleChannelName, ChannelDescriptorMarker) {
		t.Errorf("simpleChannelName is wrong, does not contain channel marker '%s', actual '%s'", ChannelDescriptorMarker, simpleChannelName)
	}
	if !strings.Contains(simpleChannelName, ChannelDescriptorSendMarker) {
		t.Errorf("simpleChannelName is wrong, does not contain channel send marker '%s', actual '%s'", ChannelDescriptorSendMarker, simpleChannelName)
	}
	if !strings.Contains(simpleChannelName, structName) {
		t.Errorf("simpleChannelName is wrong, does not contain structure  marker '%s', actual '%s'", structName, simpleChannelName)
	}
}

func TestChannelPlaceholderReplacementWithMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	channelDescriptor := NewChannelDescriptor(ChannelDirectionTx, &typeCastPlaceholder)

	switch (*(channelDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("channel reftype isn't a placeholder: '%T'", (*(channelDescriptor.RefType)))
	}

	structDescriptor := NewStructDescriptor(packageName, structName)
	var typeCastStruct TypeDescriptor = structDescriptor

	channelDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(channelDescriptor.RefType)).(type) {
	case *StructDescriptor:
	default:
		t.Errorf("replaced channel reftype isn't a struct: '%T'", (*(channelDescriptor.RefType)))
	}
}

func TestChannelPlaceholderReplacementWithNonMatchingName(t *testing.T) {
	packageName := "package"
	structName := "aStruct"

	placeholderDescriptor := NewPlaceholderDescriptor(packageName, structName)

	var typeCastPlaceholder TypeDescriptor = placeholderDescriptor

	channelDescriptor := NewChannelDescriptor(ChannelDirectionTx, &typeCastPlaceholder)

	switch (*(channelDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		t.Errorf("channel reftype isn't a placeholder: '%T'", (*(channelDescriptor.RefType)))
	}

	altStructName := "notTheStruct"
	structDescriptor := NewStructDescriptor(packageName, altStructName)
	var typeCastStruct TypeDescriptor = structDescriptor

	channelDescriptor.ReplaceRefTypeIfPlaceholder(&typeCastStruct)

	switch (*(channelDescriptor.RefType)).(type) {
	case *PlaceholderDescriptor:
	default:
		// it shouldn't have replaced because the cannonical names dont match
		t.Errorf("replaced channel reftype isn't a placeholder: '%T'", (*(channelDescriptor.RefType)))
	}
}
