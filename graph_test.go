package alice

import (
	"reflect"
	"testing"
)

func TestConstructGraph(t *testing.T) {
	var (
		rm1, _ = reflectModule(&M1{})
		rm2, _ = reflectModule(&M2{})
		rm3, _ = reflectModule(&M3{})
		rm4, _ = reflectModule(&M4{})
		rm5, _ = reflectModule(&M5{})
	)

	ms := []*reflectedModule{rm1, rm2, rm3, rm4, rm5}
	g, err := createGraph(ms...)

	if err != nil {
		t.Errorf("unexpected error after createGraph(): %s", err.Error())
	}
	if !reflect.DeepEqual(g.modules, ms) {
		t.Errorf("bad modules in graph: got %v, expected %v", g.modules, ms)
	}

	expectedG := map[*reflectedModule]map[*reflectedModule]bool{
		rm1: map[*reflectedModule]bool{
			rm2: true,
			rm4: true,
		},
		rm2: map[*reflectedModule]bool{
			rm3: true,
		},
		rm4: map[*reflectedModule]bool{
			rm2: true,
		},
		rm3: map[*reflectedModule]bool{},
		rm5: map[*reflectedModule]bool{},
	}
	if !reflect.DeepEqual(g.g, expectedG) {
		t.Errorf("bad g in graph: got %v, expected %v", g.g, expectedG)
	}
}

func TestConstructGraph_DuplicatedName(t *testing.T) {
	var (
		m1, _ = reflectModule(&M1{})
		m2, _ = reflectModule(&M1Duplicated{})
	)
	_, err := createGraph(m1, m2)

	if err == nil {
		t.Error("expect error after createGraph() of modules with duplicated name")
	}
	t.Log(err.Error())
}

func TestConstructGraph_NameNotFound(t *testing.T) {
	m, _ := reflectModule(&M4{})
	_, err := createGraph(m)

	if err == nil {
		t.Error("expect error after createGraph() of name not found")
	}
	t.Log(err.Error())
}

func TestConstructGraph_MultipleAssignableTypes(t *testing.T) {
	var (
		m1, _ = reflectModule(&M3{})
		m2, _ = reflectModule(&ModuleWithD5Impl1{})
		m3, _ = reflectModule(&ModuleWithD5Impl2{})
	)
	_, err := createGraph(m1, m2, m3)

	if err == nil {
		t.Error("expect error after createGraph() of multiple assignable types")
	}
	t.Log(err.Error())
}

func TestConstructGraph_TypeProviderNotFound(t *testing.T) {
	m, _ := reflectModule(&M3{})
	_, err := createGraph(m)

	if err == nil {
		t.Error("expect error after createGraph() of no type provider found")
	}
	t.Log(err.Error())
}

func TestConstructGraph_MultipleTypeProvider(t *testing.T) {
	var (
		m1, _ = reflectModule(&M3{})
		m2, _ = reflectModule(&ModuleWithD51{})
		m3, _ = reflectModule(&ModuleWithD52{})
	)
	_, err := createGraph(m1, m2, m3)

	if err == nil {
		t.Error("expect error after createGraph() of multiple type provider")
	}
	t.Log(err.Error())
}

func TestInstantiationOrder(t *testing.T) {
	var (
		m1, _ = reflectModule(&M1{})
		m2, _ = reflectModule(&M2{})
		m3, _ = reflectModule(&M3{})
		m4, _ = reflectModule(&M4{})
		m5, _ = reflectModule(&M5{})
	)

	ms := []*reflectedModule{m1, m2, m3, m4, m5}
	g, err := createGraph(ms...)

	if err != nil {
		t.Errorf("unexpected error after createGraph(): %s", err.Error())
	}

	expectedOrder := []*reflectedModule{m5, m1, m4, m2, m3}
	order, err := g.instantiationOrder()
	if err != nil {
		t.Errorf("unexpected error after instantiationOrder(): %s", err.Error())
	}
	if !reflect.DeepEqual(order, expectedOrder) {
		t.Errorf("bad instantiation order: got %v, expected %v", order, expectedOrder)
	}
}

func TestInstantiationOrder_Cycle(t *testing.T) {
	var (
		m1, _ = reflectModule(&M1{})
		m2, _ = reflectModule(&M2{})
		m3, _ = reflectModule(&M3{})
		m6, _ = reflectModule(&M6{})
	)

	ms := []*reflectedModule{m1, m2, m3, m6}
	g, err := createGraph(ms...)

	if err != nil {
		t.Errorf("unexpected error after createGraph(): %s", err.Error())
	}

	_, err = g.instantiationOrder()
	if err == nil {
		t.Error("expected error after instantiationOrder() with cycle")
	}
	t.Log(err.Error())
}

func TestInstantiationOrder_CycleSingleModule(t *testing.T) {
	m, _ := reflectModule(&SelfDependModule{})
	g, err := createGraph(m)

	if err != nil {
		t.Errorf("unexpected error after createGraph(): %s", err.Error())
	}

	_, err = g.instantiationOrder()
	if err == nil {
		t.Error("expected error after instantiationOrder() with single module cycle")
	}
	t.Log(err.Error())
}
