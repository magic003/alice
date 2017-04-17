package example

import (
	"reflect"
	"testing"

	"github.com/magic003/alice"
	"github.com/magic003/alice/example/business"
	"github.com/magic003/alice/example/module"
)

func TestExample(t *testing.T) {
	c := alice.CreateContainer(
		&module.ConfigModule{}, &module.PersistModule{}, &module.ClientModule{}, &module.BusinessModule{})

	retries := c.InstanceByName("Retries")
	if retries != 3 {
		t.Errorf("bad retries: got %v, expected %v", retries, 3)
	}

	table := c.InstanceByName("Table")
	if table != "example_table" {
		t.Errorf("bad table: got %s, expected %s", table, "example_table")
	}

	// should not panic
	c.Instance(reflect.TypeOf((*business.WebPageManager)(nil)))
}
