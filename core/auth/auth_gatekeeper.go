// core/auth/auth_gatekeeper.go
package auth

// Auth defines what it needs, but doesn't care WHO provides it.
type UserPrompter interface {
	RequestMFA(message string) (string, error)
}

type Authenticator struct {
	Prompter UserPrompter // This will be filled in later
}

func (a *Authenticator) Verify(user string) {
	// ... logic ...
	token, _ := a.Prompter.RequestMFA("Enter Code")
	// ... logic ...
}
