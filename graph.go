package alice

import (
	"fmt"
	"reflect"
	"strings"
)

// createGraph creates a graph of modules.
func createGraph(modules ...Module) (*graph, error) {
	g := &graph{
		modules: modules,
		g:       make(map[*reflectedModule]map[*reflectedModule]bool),
	}
	if err := g.constructGraph(); err != nil {
		return nil, err
	}
	return g, nil
}

// graph maintains the dependency relationship of the modules and gives an instantiation order.
type graph struct {
	modules []Module
	rms     []*reflectedModule
	// g is map representing the dependency graph. Modules in value depend on the key.
	// Value is a map to avoid duplication.
	g map[*reflectedModule]map[*reflectedModule]bool
}

// moduleSlice is a container of Module slice. The purpose is to be passed in recursive calls and update the slice.
type moduleSlice struct {
	modules []Module
}

// stringSlice is a container of string slice. The purpose is to be passed in recursive calls and update the slice.
type stringSlice struct {
	strings []string
}

// instantiationOrder returns the instantiation order of the modules. It returns error if there is cyclic dependencies.
func (g *graph) instantiationOrder() ([]Module, error) {
	visited := make(map[*reflectedModule]bool)
	stack := &moduleSlice{}
	recVisited := make(map[*reflectedModule]bool)
	recPath := &stringSlice{}

	for _, m := range g.rms {
		if !visited[m] {
			if err := g.dfs(m, visited, stack, recVisited, recPath); err != nil {
				return nil, err
			}
		}
	}

	return g.reverseSlice(stack.modules), nil
}

// dfs does a depth first search on the graph.
func (g *graph) dfs(
	m *reflectedModule,
	visited map[*reflectedModule]bool,
	stack *moduleSlice,
	recVisited map[*reflectedModule]bool,
	recPath *stringSlice) error {
	recPath.strings = append(recPath.strings, m.name)
	if recVisited[m] { // cyclic
		return fmt.Errorf("cyclic dependencies for modules: %s", strings.Join(recPath.strings, " -> "))
	}

	recVisited[m] = true
	for dependant := range g.g[m] {
		if !visited[dependant] {
			if err := g.dfs(dependant, visited, stack, recVisited, recPath); err != nil {
				return err
			}
		}
	}

	visited[m] = true
	stack.modules = append(stack.modules, m.m)
	recVisited[m] = false
	recPath.strings = recPath.strings[:len(recPath.strings)-1]

	return nil
}

// constructGraph constructs a graph based on the dependency of the modules.
func (g *graph) constructGraph() error {
	rmodules, nameToProviderMap, typeToProvidersMap, err := g.computeProviders()
	g.rms = rmodules
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
		if _, ok := g.g[rm]; !ok {
			g.g[rm] = make(map[*reflectedModule]bool)
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
		g.addDependencyEdge(provider, rm)
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
		g.addDependencyEdge(providers[0], rm)
	}

	return nil
}

// findAssignableProviders finds the providers which provides instances could be assigned to the specified type.
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
func (g *graph) addDependencyEdge(parent *reflectedModule, dependant *reflectedModule) {
	dependants, ok := g.g[parent]
	if !ok {
		dependants = make(map[*reflectedModule]bool)
		g.g[parent] = dependants
	}
	dependants[dependant] = true
}

// reverseSlice reverses the slice of Modules.
func (g *graph) reverseSlice(l []Module) []Module {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}

	return l
}
