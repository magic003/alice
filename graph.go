package alice

import (
	"fmt"
	"reflect"
)

func createGraph(modules ...Module) (*graph, error) {
	g := &graph{
		modules: modules,
		g:       make(map[Module]map[Module]bool),
	}
	if err := g.constructGraph(); err != nil {
		return nil, err
	}
	return g, nil
}

type graph struct {
	modules []Module
	// g is map representing the dependency graph. Modules in value depend on the key.
	// Value is a map to avoid duplication.
	g map[Module]map[Module]bool
}

func (g *graph) constructGraph() error {
	rmodules, nameToProviderMap, typeToProvidersMap, err := g.computeProviders()
	if err != nil {
		return err
	}

	// construct dependency graph
	for _, rm := range rmodules {
		if err := g.createDependenciesByNames(rm, nameToProviderMap); err != nil {
			return err
		}
		if err := g.createDependenciesByTypes(rm, typeToProvidersMap); err != nil {
			return err
		}
		if _, ok := g.g[rm.m]; !ok {
			g.g[rm.m] = make(map[Module]bool)
		}
	}

	return nil
}

// computeProviders figures out instance names and types, and the corresponding modules that provide them.
func (g *graph) computeProviders() (
	[]*reflectedModule,
	map[string]*reflectedModule,
	map[reflect.Type][]*reflectedModule,
	error) {

	var rmodules []*reflectedModule
	nameToProviderMap := make(map[string]*reflectedModule)
	typeToProvidersMap := make(map[reflect.Type][]*reflectedModule)

	for _, m := range g.modules {
		provider, err := reflectModule(m)
		if err != nil {
			return nil, nil, nil, err
		}
		rmodules = append(rmodules, provider)

		for _, name := range provider.instanceNames {
			if existingProvider, ok := nameToProviderMap[name]; ok {
				return nil, nil, nil,
					fmt.Errorf("duplicated name %s in module %s and %s", name, existingProvider.name, provider.name)
			}
			nameToProviderMap[name] = provider
		}
		for _, t := range provider.instanceTypes {
			existingProviders, _ := typeToProvidersMap[t]
			existingProviders = append(existingProviders, provider)
			typeToProvidersMap[t] = existingProviders
		}
	}

	return rmodules, nameToProviderMap, typeToProvidersMap, nil
}

// createDependenciesByNames creates dependencies of a module using its named dependencies.
func (g *graph) createDependenciesByNames(rm *reflectedModule, nameToProviderMap map[string]*reflectedModule) error {
	for _, depName := range rm.dependNames {
		provider, ok := nameToProviderMap[depName]
		if !ok {
			return fmt.Errorf("dependency name %s.%s is not found", rm.name, depName)
		}
		g.addDependencyEdge(provider.m, rm.m)
	}

	return nil
}

// createDependenciesByTypes creates dependencies of a module using its typed dependencies.
func (g *graph) createDependenciesByTypes(
	rm *reflectedModule, typeToProvidersMap map[reflect.Type][]*reflectedModule) error {
	for _, depType := range rm.dependTypes {
		providers, ok := typeToProvidersMap[depType]
		if !ok { // no exact type match, find assignable types
			assignableProviders, err := g.findAssignableProviders(rm, depType, typeToProvidersMap)
			if err != nil {
				return err
			}
			providers = assignableProviders
		}

		if len(providers) == 0 {
			return fmt.Errorf("dependency type %s.%s is not found", rm.name, depType.Name())
		}
		if len(providers) > 1 {
			var names []string
			for _, p := range providers {
				names = append(names, p.name)
			}
			return fmt.Errorf("dependency type %s.%s is found in mutiple modules: %s",
				rm.name, depType.Name(), names)
		}
		g.addDependencyEdge(providers[0].m, rm.m)
	}

	return nil
}

func (g *graph) findAssignableProviders(
	rm *reflectedModule,
	expType reflect.Type,
	typeToProvidersMap map[reflect.Type][]*reflectedModule) ([]*reflectedModule, error) {
	var providers []*reflectedModule
	foundAssignable := false
	var foundAssignableType reflect.Type
	for t, ps := range typeToProvidersMap {
		if t.AssignableTo(expType) {
			if foundAssignable {
				return nil, fmt.Errorf("multiple assignable types %s, %s for type %s.%s",
					t.Name(), foundAssignableType.Name(), rm.name, expType.Name())
			}

			providers = ps
			foundAssignable = true
			foundAssignableType = t
		}
	}

	return providers, nil
}

// addDependencyEdge creates a dependency edge in the graph. dependant depends on parent.
func (g *graph) addDependencyEdge(parent Module, dependant Module) {
	dependants, ok := g.g[parent]
	if !ok {
		dependants = make(map[Module]bool)
		g.g[parent] = dependants
	}
	dependants[dependant] = true
}
