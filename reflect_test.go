package alice

import (
	"reflect"
	"testing"
)

type Dep1 interface {
	IsDep1() bool
}

type Dep1Impl struct{}

func (d *Dep1Impl) IsDep1() bool {
	return true
}

type Dep2 interface {
	IsDep2() bool
}

type Dep2Impl struct{}

func (d *Dep2Impl) IsDep2() bool {
	return true
}

type testModule struct {
	BaseModule
	dep1   Dep1 `alice:""`
	dep2   Dep2 `alice:"Dep2"`
	nonDep string
}

func (m *testModule) Dep1() Dep1 {
	return &Dep1Impl{}
}

func (m *testModule) Dep2() Dep2 {
	return &Dep2Impl{}
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

func (m *invalidMethodModule1) Dep1(str string) Dep1 {
	return &Dep1Impl{}
}

type invalidMethodModule2 struct {
	BaseModule
}

func (m *invalidMethodModule2) Dep2() (Dep2, error) {
	return &Dep2Impl{}, nil
}

func TestReflectModule(t *testing.T) {
	m := &testModule{}

	rmodule, err := reflectModule(m)

	if err != nil {
		t.Errorf("unexpected error after reflectModule(): %s", err.Error())
	}
	if rmodule.m != m {
		t.Errorf("bad m in reflectedModule: got %v, expected %v", rmodule.m, m)
	}
	expectedName := "testModule"
	if rmodule.name != expectedName {
		t.Errorf("bad name in reflectedModule: got %s, expected %s", rmodule.name, expectedName)
	}

	expectedInstanceNames := []string{"Dep1", "Dep2"}
	if !reflect.DeepEqual(rmodule.instanceNames, expectedInstanceNames) {
		t.Errorf("bad instanceNames in reflectedModule: got %v, expected %v",
			rmodule.instanceNames, expectedInstanceNames)
	}
	expectedInstanceTypes := []reflect.Type{reflect.TypeOf((*Dep1)(nil)).Elem(), reflect.TypeOf((*Dep2)(nil)).Elem()}
	if !reflect.DeepEqual(rmodule.instanceTypes, expectedInstanceTypes) {
		t.Errorf("bad instanceTypes in reflectedModule: got %v, expected %v",
			rmodule.instanceTypes, expectedInstanceTypes)
	}

	expectedDependNames := []string{"Dep2"}
	if !reflect.DeepEqual(rmodule.dependNames, expectedDependNames) {
		t.Errorf("bad dependNames in reflectedModule: got %v, expected %v",
			rmodule.dependNames, expectedDependNames)
	}
	expectedDependTypes := []reflect.Type{reflect.TypeOf((*Dep1)(nil)).Elem()}
	if !reflect.DeepEqual(rmodule.dependTypes, expectedDependTypes) {
		t.Errorf("bad dependTypes in reflectedModule: got %v, expected %v",
			rmodule.dependTypes, expectedDependTypes)
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
