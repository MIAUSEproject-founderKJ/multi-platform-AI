//runtime/interface_adapter/hmi_adapter.go - implementation layer only

package interface_adapter

import (
	"fmt"

	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

// ===================Screen Adapter
type ScreenAdapter struct{}

func NewScreenAdapter() *ScreenAdapter {
	return &ScreenAdapter{}
}

func (s *ScreenAdapter) Start(session *user_setting.UserSession) error {
	fmt.Println("Screen adapter started")
	return nil
}

func (s *ScreenAdapter) Notify(msg string) {
	fmt.Println("[SCREEN]", msg)
}
