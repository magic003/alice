package alice

import (
	"reflect"
	"testing"
)

type reflectTestModule struct {
	BaseModule
	dep1   D1 `alice:""`
	dep2   D2 `alice:"Dep2"`
	nonDep string
}

func (m *reflectTestModule) Dep1() D1 {
	return &D1Impl{}
}

func (m *reflectTestModule) Dep2() D2 {
	return &D2Impl{}
}

type nonPointerModule struct{}

func (m nonPointerModule) IsModule() bool {
	return true
}

type nonStructModule string

func (m *nonStructModule) IsModule() bool {
	return true
}

type invalidMethodModule1 struct {
	BaseModule
}

func (m *invalidMethodModule1) Dep1(str string) D1 {
	return &D1Impl{}
}

type invalidMethodModule2 struct {
	BaseModule
}

func (m *invalidMethodModule2) Dep2() (D2, error) {
	return &D2Impl{}, nil
}

func TestReflectModule(t *testing.T) {
	m := &reflectTestModule{}

	rmodule, err := reflectModule(m)

	if err != nil {
		t.Errorf("unexpected error after reflectModule(): %s", err.Error())
	}
	if rmodule.m != m {
		t.Errorf("bad m in reflectedModule: got %v, expected %v", rmodule.m, m)
	}
	expectedName := "reflectTestModule"
	if rmodule.name != expectedName {
		t.Errorf("bad name in reflectedModule: got %s, expected %s", rmodule.name, expectedName)
	}

	expectedInstances := []*instanceMethod{
		{
			name:   "Dep1",
			tp:     reflect.TypeOf((*D1)(nil)).Elem(),
			method: reflect.ValueOf(m).MethodByName("Dep1"),
		},
		{
			name:   "Dep2",
			tp:     reflect.TypeOf((*D2)(nil)).Elem(),
			method: reflect.ValueOf(m).MethodByName("Dep2"),
		},
	}
	if !reflect.DeepEqual(rmodule.instances, expectedInstances) {
		t.Errorf("bad instances in reflectedModule: got %v, expected %v",
			rmodule.instances, expectedInstances)
	}

	expectedNamedDepends := []*namedField{
		{
			name:  "Dep2",
			field: reflect.ValueOf(m).Elem().FieldByName("dep2"),
		},
	}
	if !reflect.DeepEqual(rmodule.namedDepends, expectedNamedDepends) {
		t.Errorf("bad namedDepends in reflectedModule: got %v, expected %v",
			rmodule.namedDepends, expectedNamedDepends)
	}

	expectedTypedDpends := []*typedField{
		{
			tp:    reflect.TypeOf((*D1)(nil)).Elem(),
			field: reflect.ValueOf(m).Elem().FieldByName("dep1"),
		},
	}
	if !reflect.DeepEqual(rmodule.typedDepends, expectedTypedDpends) {
		t.Errorf("bad typedDepends in reflectedModule: got %v, expected %v",
			rmodule.typedDepends, expectedTypedDpends)
	}
}

func TestReflectModule_InvalidModuleType(t *testing.T) {
	nonPtrModule := nonPointerModule{}

	_, err := reflectModule(nonPtrModule)
	if err == nil {
		t.Error("expect error after reflectModule() on non-pointer module")
	}
	t.Log(err.Error())

	nonStructModule := nonStructModule("module")
	_, err = reflectModule(&nonStructModule)
	if err == nil {
		t.Error("expect error after reflectModule() on non-struct module")
	}
	t.Log(err.Error())
}

func TestReflectModule_InvalidMethod(t *testing.T) {
	m1 := &invalidMethodModule1{}

	_, err := reflectModule(m1)
	if err == nil {
		t.Error("expect error after reflectModule() on module with 1 paramter method")
	}
	t.Log(err.Error())

	m2 := &invalidMethodModule2{}
	_, err = reflectModule(m2)
	if err == nil {
		t.Error("expect error after reflectModule() on module with 2 return values method")
	}
	t.Log(err.Error())
}
