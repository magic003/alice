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

	instances    []*instanceMethod
	namedDepends []*namedField
	typedDepends []*typedField
}

type instanceMethod struct {
	name   string
	tp     reflect.Type
	method reflect.Value
}

type namedField struct {
	name  string
	field reflect.Value
}

type typedField struct {
	tp    reflect.Type
	field reflect.Value
}

// reflectModule creates a reflectedModule from a Module. It returns error if the Module is not properly defined.
func reflectModule(m Module) (*reflectedModule, error) {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("module %s is not a pointer of struct", v.String())
	}

	// get instances
	ptrT := v.Type()
	var instances []*instanceMethod
	for i := 0; i < ptrT.NumMethod(); i++ {
		method := ptrT.Method(i)
		if method.Name == _IsModuleMethodName {
			continue
		}
		if method.Type.NumIn() != 1 || method.Type.NumOut() != 1 { // receiver is the first parameter
			return nil, fmt.Errorf("method %s.%s doesn't have 0 parameter and 1 return value",
				v.Elem().Type().Name(), method.Name)
		}
		instances = append(instances, &instanceMethod{
			name:   method.Name,
			tp:     method.Type.Out(0),
			method: v.MethodByName(method.Name),
		})
	}

	// get dependencies
	t := v.Elem().Type()
	var namedDepends []*namedField
	var typedDepends []*typedField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		if dependName, exists := field.Tag.Lookup(_Tag); exists {
			if dependName != "" {
				namedDepends = append(namedDepends, &namedField{
					name:  dependName,
					field: v.Elem().FieldByName(field.Name),
				})
			} else {
				typedDepends = append(typedDepends, &typedField{
					tp:    field.Type,
					field: v.Elem().FieldByName(field.Name),
				})
			}
		}
	}

	return &reflectedModule{
		m:            m,
		name:         t.Name(),
		instances:    instances,
		namedDepends: namedDepends,
		typedDepends: typedDepends,
	}, nil
}
