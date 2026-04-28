//internal/policy/decision.go
package policy

type Decision struct {
	Allowed bool
	Reason  string

	Granted []string
	Denied  []string
}
