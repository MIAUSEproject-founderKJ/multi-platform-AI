// modules/kernel_extension/adapters/legacy/adapter.go
package module_legacy

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/data_transport/ingress"
	"github.com/tetratelabs/wazero/api"
)

type WASMModule struct {
	module api.Module
}

func (w *WASMModule) Name() string {
	return w.module.Name()
}

func (w *WASMModule) Handle(ctx context.Context, payload []byte) error {
	fn := w.module.ExportedFunction("handle")

	mem := w.module.Memory()

	ptr, ok := mem.Allocate(uint32(len(payload)))
	if !ok {
		return fmt.Errorf("memory allocation failed")
	}

	mem.Write(ptr, payload)

	_, err := fn.Call(ctx, uint64(ptr), uint64(len(payload)))
	return err
}

var _ ingress.Handler = (*WASMModule)(nil)
