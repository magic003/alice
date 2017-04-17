package module

import (
	"github.com/magic003/alice"
)

// ConfigModule is the module for configurations.
type ConfigModule struct {
	alice.BaseModule
}

// Retries returns the client retry attempts.
func (m *ConfigModule) Retries() int {
	return 3
}

// Table is the table name.
func (m *ConfigModule) Table() string {
	return "example_table"
}
