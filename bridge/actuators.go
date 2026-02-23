// bridge/actuators.go
package bridge

import (
	"fmt"
)

type SafeBridge struct {
	Inner  PowerController
	Kernel *Kernel
}

func (s *SafeBridge) WriteActuator(name string, value float64) error {
	if !s.Kernel.CanActuate() {
		return fmt.Errorf("actuation blocked: insufficient trust or mode")
	}
	return s.Inner.WriteActuator(name, value)
}
