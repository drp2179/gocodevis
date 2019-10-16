package types

import "strings"

const (
	//ChannelDescriptorMarker - the marker for a channel
	ChannelDescriptorMarker = "chan"
	// ChannelDescriptorSendMarker - the marker for a transmit channel
	ChannelDescriptorSendMarker = "->"
	// ChannelDescriptorReceiveMarker - the marker for a receive channel
	ChannelDescriptorReceiveMarker = "<-"
	// ChannelDescriptorNoneMarker - marker for a unqualified channel
	ChannelDescriptorNoneMarker = "--"

	// ChannelDirectionNone - unqualified channel
	ChannelDirectionNone = 0
	// ChannelDirectionRx -  receive channel
	ChannelDirectionRx = 1
	// ChannelDirectionTx - transmit channel
	ChannelDirectionTx = 2
)

// ChannelDescriptor - represents a channel type
type ChannelDescriptor struct {
	Direction int
	RefType   *TypeDescriptor
}

// NewChannelDescriptor = create a new channel descriptor
func NewChannelDescriptor(channelDirection int, refType *TypeDescriptor) *ChannelDescriptor {
	ptr := ChannelDescriptor{}
	ptr.RefType = refType
	ptr.Direction = channelDirection

	return &ptr
}

// GetChannelTypeNameForTypename - get the channel type name for the given type name and direction
func GetChannelTypeNameForTypename(channelDirection int, typeName string) string {
	if channelDirection == ChannelDirectionRx {
		return ChannelDescriptorMarker + ChannelDescriptorReceiveMarker + typeName
	} else if channelDirection == ChannelDirectionTx {
		return ChannelDescriptorMarker + ChannelDescriptorSendMarker + typeName
	}
	return ChannelDescriptorMarker + ChannelDescriptorNoneMarker + typeName
}

func (channel ChannelDescriptor) String() string {
	sb := strings.Builder{}

	sb.WriteString(channel.GetCannonicalName())

	if channel.Direction == ChannelDirectionRx {
		sb.WriteString(" is receive channel of ")
	} else if channel.Direction == ChannelDirectionTx {
		sb.WriteString(" is transmit channel of ")
	} else {
		sb.WriteString(" is channel of ")
	}

	sb.WriteString((*(channel.RefType)).GetCannonicalName())
	sb.WriteString("\n")

	return sb.String()
}

// --------------------------
// Type interface
//

// GetSimpleName - Type interface to get a name
func (channel ChannelDescriptor) GetSimpleName() string {
	return GetChannelTypeNameForTypename(channel.Direction, (*(channel.RefType)).GetSimpleName())
}

// GetCannonicalName - Type interface to get the cannonical name
func (channel ChannelDescriptor) GetCannonicalName() string {
	return GetChannelTypeNameForTypename(channel.Direction, (*(channel.RefType)).GetCannonicalName())
}

// ReplaceRefTypeIfPlaceholder - replace the pointer type
func (channel *ChannelDescriptor) ReplaceRefTypeIfPlaceholder(refType *TypeDescriptor) {
	switch (*(channel.RefType)).(type) {
	case *PlaceholderDescriptor:
		channelRefTypeName := (*(channel.RefType)).GetCannonicalName()
		newRefTypeName := (*(refType)).GetCannonicalName()

		if channelRefTypeName == newRefTypeName {
			channel.RefType = refType
		}
	}
}
