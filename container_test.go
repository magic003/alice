package alice

import (
	"reflect"
	"testing"
)

//***********************************************************
// Definitions of dependencies and modules used for testing
//***********************************************************

type D1 interface {
	D1()
}

type D1Impl struct{}

func (d *D1Impl) D1() {}

type D2 interface {
	D2()
}

type D2Impl struct{}

func (d *D2Impl) D2() {}

type D3 interface {
	D3()
}

type D3Impl struct{}

func (d *D3Impl) D3() {}

type D4 interface {
	D4()
}

type D4Impl struct{}

func (d *D4Impl) D4() {}

type D5 interface {
	D5()
}

type D5Impl struct{}

func (d *D5Impl) D5() {}

type D5Impl2 struct{}

func (d *D5Impl2) D5() {}

type M1 struct {
	BaseModule
}

func (m *M1) D1() D1 {
	return &D1Impl{}
}

func (m *M1) D2() D2 {
	return &D2Impl{}
}

type M2 struct {
	BaseModule
	D1 D1 `alice:"D1"`
	D2 D2 `alice:"D2"`
	D3 D3 `alice:""`
	D4 D4 `alice:""`
}

func (m *M2) D5() *D5Impl {
	return &D5Impl{}
}

type M3 struct {
	BaseModule
	D5 D5 `alice:""`
}

func (m *M3) DM3() D1 {
	return &D1Impl{}
}

type M4 struct {
	BaseModule
	D1 D1 `alice:"D1"`
}

func (m *M4) D3() D3 {
	return &D3Impl{}
}

func (m *M4) D4() D4 {
	return &D4Impl{}
}

type M5 struct {
	BaseModule
}

type M6 struct {
	BaseModule
	D1 D1 `alice:"DM3"`
}

func (m *M6) D3() D3 {
	return &D3Impl{}
}

func (m *M6) D4() D4 {
	return &D4Impl{}
}

type M1Duplicated struct {
	BaseModule
}

func (m *M1Duplicated) D1() D1 {
	return &D1Impl{}
}

type ModuleWithD51 struct {
	BaseModule
}

func (m *ModuleWithD51) D5_1() D5 {
	return &D5Impl{}
}

type ModuleWithD52 struct {
	BaseModule
}

func (m *ModuleWithD52) D5_2() D5 {
	return &D5Impl2{}
}

type ModuleWithD5Impl1 struct {
	BaseModule
}

func (m *ModuleWithD5Impl1) D5_1() *D5Impl {
	return &D5Impl{}
}

type ModuleWithD5Impl2 struct {
	BaseModule
}

func (m *ModuleWithD5Impl2) D5_2() *D5Impl2 {
	return &D5Impl2{}
}

type SelfDependModule struct {
	BaseModule
	D D1 `alice:"D1"`
}

func (m *SelfDependModule) D1() D1 {
	return &D1Impl{}
}

//***********************************************************

func TestPopulate(t *testing.T) {
	var (
		m1 = &M1{}
		m2 = &M2{}
		m3 = &M3{}
		m4 = &M4{}
		m5 = &M5{}
	)

	c := &container{modules: []Module{m1, m2, m3, m4, m5}}
	c.populate()

	expectedM2 := &M2{
		D1: &D1Impl{},
		D2: &D2Impl{},
		D3: &D3Impl{},
		D4: &D4Impl{},
	}
	if !reflect.DeepEqual(m2, expectedM2) {
		t.Errorf("bad m2 after populate(): got %v, expected %v", m2, expectedM2)
	}
	expectedM3 := &M3{
		D5: &D5Impl{},
	}
	if !reflect.DeepEqual(m3, expectedM3) {
		t.Errorf("bad m3 after populate(): got %v, expected %v", m3, expectedM3)
	}
	expectedM4 := &M4{
		D1: &D1Impl{},
	}
	if !reflect.DeepEqual(m4, expectedM4) {
		t.Errorf("bad m4 after populate(): got %v, expected %v", m4, expectedM4)
	}
	expectedInstanceByName := map[string]interface{}{
		"D1":  &D1Impl{},
		"D2":  &D2Impl{},
		"D5":  &D5Impl{},
		"DM3": &D1Impl{},
		"D3":  &D3Impl{},
		"D4":  &D4Impl{},
	}
	if !reflect.DeepEqual(c.instanceByName, expectedInstanceByName) {
		t.Errorf("bad instanceByName after populate(): got %v, expected %v", c.instanceByName, expectedInstanceByName)
	}

	expectedInstanceByType := map[reflect.Type][]interface{}{
		reflect.TypeOf((*D1)(nil)).Elem(): []interface{}{
			&D1Impl{},
			&D1Impl{},
		},
		reflect.TypeOf((*D2)(nil)).Elem(): []interface{}{
			&D2Impl{},
		},
		reflect.TypeOf((*D5Impl)(nil)): []interface{}{
			&D5Impl{},
		},
		reflect.TypeOf((*D3)(nil)).Elem(): []interface{}{
			&D3Impl{},
		},
		reflect.TypeOf((*D4)(nil)).Elem(): []interface{}{
			&D4Impl{},
		},
	}
	if !reflect.DeepEqual(c.instanceByType, expectedInstanceByType) {
		t.Errorf("bad instanceByType after populate(): got %v, expected %v", c.instanceByType, expectedInstanceByType)
	}
}

func TestPopulate_PanicOnInvalidModule(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic after populate() on invalid module")
		} else {
			t.Log(r)
		}
	}()
	c := &container{modules: []Module{nonPointerModule{}}}
	c.populate()
}

func TestPopulate_PanicOnCreateGraphError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic after populate() on create graph error")
		} else {
			t.Log(r)
		}
	}()
	c := &container{modules: []Module{&M1{}, &M1Duplicated{}}}
	c.populate()
}

func TestPopulate_PanicOnInstantiationOrderError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic after populate() on instantiation order error")
		} else {
			t.Log(r)
		}
	}()
	c := &container{modules: []Module{&M1{}, &M2{}, &M3{}, &M6{}}}
	c.populate()
}

func TestInstance(t *testing.T) {
	var (
		m1 = &M1{}
		m2 = &M2{}
		m3 = &M3{}
		m4 = &M4{}
		m5 = &M5{}
	)

	c := &container{modules: []Module{m1, m2, m3, m4, m5}}
	c.populate()

	d2 := c.Instance(reflect.TypeOf((*D2)(nil)).Elem()).(D2)
	expectedD2 := &D2Impl{}
	if !reflect.DeepEqual(d2, expectedD2) {
		t.Errorf("bad instance from Instance(): got %v, expected %v", d2, expectedD2)
	}

	d5 := c.Instance(reflect.TypeOf((*D5Impl)(nil))).(*D5Impl)
	expectedD5 := &D5Impl{}
	if !reflect.DeepEqual(d5, expectedD5) {
		t.Errorf("bad instance from Instance(): got %v, expected %v", d5, expectedD5)
	}
}

func TestInstance_PanicOnTypeNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for Instance() on type not found")
		} else {
			t.Log(r)
		}
	}()

	c := &container{modules: []Module{&M1{}}}
	c.populate()

	c.Instance(reflect.TypeOf((*D3)(nil)).Elem())
}

func TestInstance_PanicOnMultipleMatchedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for Instance() on multiple matched type")
		} else {
			t.Log(r)
		}
	}()

	c := &container{modules: []Module{&M1{}, &M2{}, &M3{}, &M4{}}}
	c.populate()

	c.Instance(reflect.TypeOf((*D1)(nil)).Elem())
}

func TestInstanceByName(t *testing.T) {
	var (
		m1 = &M1{}
		m2 = &M2{}
		m3 = &M3{}
		m4 = &M4{}
		m5 = &M5{}
	)

	c := &container{modules: []Module{m1, m2, m3, m4, m5}}
	c.populate()

	d1 := c.InstanceByName("D1").(D1)
	expectedD1 := &D1Impl{}
	if !reflect.DeepEqual(d1, expectedD1) {
		t.Errorf("bad instance from InstanceByName(): got %v, expected %v", d1, expectedD1)
	}

	dm3 := c.InstanceByName("DM3").(D1)
	expectedDM3 := &D1Impl{}
	if !reflect.DeepEqual(dm3, expectedDM3) {
		t.Errorf("bad instance from InstanceByName(): got %v, expected %v", dm3, expectedDM3)
	}
}

func TestInstanceByName_PanicOnNameNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for InstanceByName() on name not found")
		} else {
			t.Log(r)
		}
	}()
	c := &container{modules: []Module{&M1{}}}
	c.populate()

	c.InstanceByName("D3")
}

func TestCreateContainer(t *testing.T) {
	var (
		m1 = &M1{}
		m2 = &M2{}
		m3 = &M3{}
		m4 = &M4{}
		m5 = &M5{}
	)

	c := CreateContainer(m1, m2, m3, m4, m5)

	d1 := c.InstanceByName("D1").(D1)
	expectedD1 := &D1Impl{}
	if !reflect.DeepEqual(d1, expectedD1) {
		t.Errorf("bad instance after CreateContainer(): got %v, expected %v", d1, expectedD1)
	}
}
