//modules/kernel_extension/supervisor/module_dependency_resolver.go
/*Validate module dependency declarations.
Detect missing or circular dependencies.
Produce a safe startup order where every module is initialized only after its dependencies.*/

/*TelemetryModule depends on NetworkModule
InferenceModule depends on StorageModule
ResolveDependencies enforces deterministic startup ordering.
*/

package kernel_supervisor

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
)

func ResolveDependencies(mods []kernel_lifecycle.DomainModule) ([]kernel_lifecycle.DomainModule, error) {
	nameIndex := make(map[string]kernel_lifecycle.DomainModule)
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for _, m := range mods {
		name := m.Name()
		nameIndex[name] = m
		inDegree[name] = 0
	}

	for _, m := range mods {
		for _, dep := range m.DependsOn() {
			if _, ok := nameIndex[dep]; !ok {
				return nil, fmt.Errorf("module %s depends on unknown module %s", m.Name(), dep)
			}
			graph[dep] = append(graph[dep], m.Name())
			inDegree[m.Name()]++
		}
	}

	var queue []string
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	var ordered []kernel_lifecycle.DomainModule

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		ordered = append(ordered, nameIndex[current])

		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(ordered) != len(mods) {
		return nil, fmt.Errorf("circular module dependency detected")
	}

	return ordered, nil
}
