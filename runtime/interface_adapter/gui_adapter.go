// runtime/interface_adapter/gui_adapter.go
package interface_adapter

import (
	auth "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"

	"fmt"
)

type GUIAdapter struct{}

func (g *GUIAdapter) Start(session *user_setting.UserSession) error {
	// Launch window
	return nil
}

func (c *GUIAdapter) Notify(msg string) {
	fmt.Println("[GUI]", msg)
}
func (g *GUIAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

type GUIAuth struct{}

func NewGUIAuth() auth.AuthInterface {
	return &GUIAuth{}
}
func (g *GUIAuth) Authenticate() error {
	return nil
}
