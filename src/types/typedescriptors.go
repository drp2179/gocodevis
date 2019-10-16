package types

import (
	"strings"
)

// TypeDescriptor - defines a generic type
type TypeDescriptor interface {
	GetSimpleName() string
	GetCannonicalName() string
}

// IsNameLocalToScope - given the name and the scope, return the resolved name, any extracged prefixes and if its local to the scope
func IsNameLocalToScope(name string, scopeName string) (resolvedName string, extractedPrefix string, isLocal bool) {

	if len(scopeName) > 0 && strings.HasPrefix(name, scopeName) {
		localName := name[len(scopeName)+1:]
		return localName, "", true
	}

	if strings.HasPrefix(name, ArrayDescriptorMarker) {
		newName := name[len(ArrayDescriptorMarker):]
		resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(newName, scopeName)
		return resolvedName, ArrayDescriptorMarker + extractedPrefix, isLocal
	} else if strings.HasPrefix(name, EllipsisDescriptorMarker) {
		newName := name[len(EllipsisDescriptorMarker):]
		resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(newName, scopeName)
		return resolvedName, EllipsisDescriptorMarker + extractedPrefix, isLocal
	} else if strings.HasPrefix(name, PointerDescriptorMarker) {
		newName := name[len(PointerDescriptorMarker):]
		resolvedName, extractedPrefix, isLocal := IsNameLocalToScope(newName, scopeName)
		return resolvedName, PointerDescriptorMarker + extractedPrefix, isLocal
	}

	// TODO: do we need Map or Channel support here?

	return name, "", false
}
