package alice

import (
	"reflect"
)

// Container defines the interface of an instance container. It initializes instances based on dependencies,
// and provides APIs to retrieve instances by type or name.
type Container interface {
	// Instance returns an instance by type. It panics if there are multiple instances for the same type.
	// It returns nil if no instance is found.
	Instance(t reflect.Type) interface{}
	// InstanceByName returns an instance by name. It returns nil if no instance is found.
	InstanceByName(name string) interface{}
}
