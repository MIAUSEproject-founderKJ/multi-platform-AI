//runtime/interface_adapter/tui_adapter.go

package interface_adapter

import (
	"fmt"

	auth "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user_setting"
)

type TUIAdapter struct{}

func (t *TUIAdapter) Start(session *user_setting.UserSession) error {
	// integrate charmbracelet/bubbletea
	return nil
}
func (c *TUIAdapter) Notify(msg string) {
	fmt.Println("[TUI]", msg)
}

func (t *TUIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

type TUIAuth struct{}

func NewTUIAuth() auth.AuthInterface {
	return &TUIAuth{}
}

func (t *TUIAuth) Authenticate() error {
	return nil
}
