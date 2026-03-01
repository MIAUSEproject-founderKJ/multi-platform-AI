//modules/implement_unknown/wasm_loader.go

func LoadWASMModule(path string) (*WASMModule, error) {

    ctx := context.Background()
    r := wazero.NewRuntime(ctx)

    wasmBytes, _ := os.ReadFile(path)

    mod, err := r.InstantiateModuleFromBinary(ctx, wasmBytes)
    if err != nil {
        return nil, err
    }

    return &WASMModule{module: mod}, nil
}