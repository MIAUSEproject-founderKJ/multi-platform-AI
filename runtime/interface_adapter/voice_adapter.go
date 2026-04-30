//runtime/interface_adapter/voice_adapter.go

package interface_adapter

func (v *VoiceAuth) StartAuthFlow(am *auth.AuthManager) (*user_setting.UserSession, error) {
	return am.LoginOrSignUpInteractive()
}

type VoiceAuth struct{}

func NewVoiceAuth() auth.AuthInterface {
	return &VoiceAuth{}
}
func (v *VoiceAuth) Authenticate() error {
	return nil
}
