//main/main.go

package main

import (
	"aios/core"
	"aios/modules"
)

func main() {

	ctx := ResolveRuntimeContext() // hardware probe + identity + boot state

	registry := core.NewRegistry()
	registry.Register(&modules.AutonomousKernel{})
	registry.Register(&modules.ProductivityModule{})

	loader := core.NewLoader(registry)

	if err := loader.ResolveAndLoad(ctx); err != nil {
		panic(err)
	}

	if err := loader.StartAll(); err != nil {
		panic(err)
	}

	select {} // block
}
