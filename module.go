package alice

// Module is a marker interface for structs that defines how to initialize instances.
type Module interface {
	// IsModule indicates if this is a module.
	IsModule() bool
}

// BaseModule is an implementation of Module interface. It should be embeded into each module defined in the
// application.
//
// A typical module is defined as follows:
//
//	type ExampleModule struct {
//		alice.BaseModule
//		Foo Foo `alice:""`		// associated by type
//		Bar Bar `alice:"Bar"`	// associated by name
//		URL string				// not associated. Provided by creating the module.
//	}
//
//	func (m *ExampleModule) Baz() Baz {
//		return Baz{}
//	}
type BaseModule struct{}

// IsModule indicates it is a module.
func (b *BaseModule) IsModule() bool {
	return true
}
