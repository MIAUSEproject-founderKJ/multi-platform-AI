// core/platform/probe/heartbeat.go

package probe

import (
	"errors"
	"fmt"
)

func Heartbeat(currentID *HardwareIdentity, savedProfile *HardwareProfile) error {
	// 1. Check if the physical machine changed
	if currentID.InstanceID != savedProfile.ID {
		return errors.New("instance_mismatch")
	}

	// 2. Check if a critical sensor disappeared
	for _, sensor := range savedProfile.ActiveSensors {
		if !isSensorAlive(sensor) {
			return fmt.Errorf("missing_critical_sensor: %s", sensor)
		}
	}
	return nil
}
