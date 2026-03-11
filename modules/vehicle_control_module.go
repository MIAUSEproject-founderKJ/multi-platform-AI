//modules/vehicle_control_module.go

package modules

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
)

type VehicleControlModule struct {
	ctx     *runtime.RuntimeContext
	healthy atomic.Bool
}

func NewVehicleControlModule() DomainModule {
	return &VehicleControlModule{}
}

func (m *VehicleControlModule) Name() string {
	return "VehicleControlModule"
}

func (m *VehicleControlModule) DependsOn() []string {
	return []string{"InferenceModule"}
}

func (m *VehicleControlModule) Init(ctx *runtime.RuntimeContext) error {
	m.ctx = ctx
	m.healthy.Store(true)
	return nil
}

func (m *VehicleControlModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *VehicleControlModule) Handle(ctx context.Context, payload []byte) error {

	if len(payload) == 0 {
		return errors.New("empty control command")
	}

	m.ctx.Logger.Info("vehicle command executed",
		"cmd", string(payload),
	)

	return nil
}

func (m *VehicleControlModule) Healthy() bool {
	return m.healthy.Load()
}