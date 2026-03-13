// modules/implement_unknown/adapter.go
package implement_unknown

import (
	"context"
	"unsafe"

	"github.com/tetratelabs/wazero/api"
)

type WASMModule struct {
	module api.Module
}

func (w *WASMModule) Name() string { return "DynamicModule" }

func (w *WASMModule) Handle(ctx context.Context, payload []byte) error {
	fn := w.module.ExportedFunction("handle")
	_, err := fn.Call(ctx, uint64(uintptr(unsafe.Pointer(&payload[0]))))
	return err
}
