package alice

// Module is a marker interface for structs that defines how to initialize instances.
type Module interface {
	// IsModule indicates if this is a module.
	IsModule() bool
}

// BaseModule is an implementation of Module interface. It should be embeded into each module defined in the
// application.
type BaseModule struct{}

// IsModule indicates it is a module.
func (b *BaseModule) IsModule() bool {
	return true
}
