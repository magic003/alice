package alice

import (
	"fmt"
	"reflect"
)

const _Tag = "alice"
const _IsModuleMethodName = "IsModule"

// reflectedModule contains the instance and dependency information of a Module. The information is extracted
// using reflection.
type reflectedModule struct {
	m    Module
	name string

	instanceNames []string
	instanceTypes []reflect.Type
	dependNames   []string
	dependTypes   []reflect.Type
}

// reflectModule creates a reflectedModule from a Module. It returns error if the Module is not properly defined.
func reflectModule(m Module) (*reflectedModule, error) {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("module %s is not a pointer of struct", v.String())
	}

	// get instances
	ptrT := v.Type()
	var instanceNames []string
	var instanceTypes []reflect.Type
	for i := 0; i < ptrT.NumMethod(); i++ {
		method := ptrT.Method(i)
		if method.Name == _IsModuleMethodName {
			continue
		}
		if method.Type.NumIn() != 1 || method.Type.NumOut() != 1 { // receiver is the first parameter
			return nil, fmt.Errorf("method %s.%s doesn't have 0 parameter and 1 return value",
				v.Elem().Type().Name(), method.Name)
		}
		instanceNames = append(instanceNames, method.Name)
		instanceTypes = append(instanceTypes, method.Type.Out(0))
	}

	// get dependencies
	t := v.Elem().Type()
	var dependNames []string
	var dependTypes []reflect.Type
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		if dependName, exists := field.Tag.Lookup(_Tag); exists {
			if dependName != "" {
				dependNames = append(dependNames, dependName)
			} else {
				dependTypes = append(dependTypes, field.Type)
			}
		}
	}

	return &reflectedModule{
		m:             m,
		name:          t.Name(),
		instanceNames: instanceNames,
		instanceTypes: instanceTypes,
		dependNames:   dependNames,
		dependTypes:   dependTypes,
	}, nil
}
