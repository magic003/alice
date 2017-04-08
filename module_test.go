package alice

import (
	"testing"
)

func TestBaseModule(t *testing.T) {
	var base Module = &BaseModule{}
	if !base.IsModule() {
		t.Error("BaseModule is expected to be a Module, but it is not")
	}
}
