package alice

import (
	"reflect"
	"testing"
)

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

func TestConstructGraph(t *testing.T) {
	var (
		m1 = &M1{}
		m2 = &M2{}
		m3 = &M3{}
		m4 = &M4{}
		m5 = &M5{}
	)

	ms := []Module{m1, m2, m3, m4, m5}
	g, err := createGraph(ms...)

	if err != nil {
		t.Errorf("unexpected error after createGraph(): %s", err.Error())
	}
	if !reflect.DeepEqual(g.modules, ms) {
		t.Errorf("bad modules in graph: got %v, expected %v", g.modules, ms)
	}

	expectedG := map[Module]map[Module]bool{
		m1: map[Module]bool{
			m2: true,
			m4: true,
		},
		m2: map[Module]bool{
			m3: true,
		},
		m4: map[Module]bool{
			m2: true,
		},
		m3: map[Module]bool{},
		m5: map[Module]bool{},
	}
	if !reflect.DeepEqual(g.g, expectedG) {
		t.Errorf("bad g in graph: got %v, expected %v", g.g, expectedG)
	}
}

func TestConstructGraph_InvalidModule(t *testing.T) {
	_, err := createGraph(&invalidMethodModule1{})

	if err == nil {
		t.Error("expect error after createGraph() of invalid module")
	}
	t.Log(err.Error())
}

func TestConstructGraph_DuplicatedName(t *testing.T) {
	_, err := createGraph(&M1{}, &M1Duplicated{})

	if err == nil {
		t.Error("expect error after createGraph() of modules with duplicated name")
	}
	t.Log(err.Error())
}

func TestConstructGraph_NameNotFound(t *testing.T) {
	_, err := createGraph(&M4{})

	if err == nil {
		t.Error("expect error after createGraph() of name not found")
	}
	t.Log(err.Error())
}

func TestConstructGraph_MultipleAssignableTypes(t *testing.T) {
	_, err := createGraph(&M3{}, &ModuleWithD5Impl1{}, &ModuleWithD5Impl2{})

	if err == nil {
		t.Error("expect error after createGraph() of multiple assignable types")
	}
	t.Log(err.Error())
}

func TestConstructGraph_TypeProviderNotFound(t *testing.T) {
	_, err := createGraph(&M3{})

	if err == nil {
		t.Error("expect error after createGraph() of no type provider found")
	}
	t.Log(err.Error())
}

func TestConstructGraph_MultipleTypeProvider(t *testing.T) {
	_, err := createGraph(&M3{}, &ModuleWithD51{}, &ModuleWithD52{})

	if err == nil {
		t.Error("expect error after createGraph() of multiple type provider")
	}
	t.Log(err.Error())
}
