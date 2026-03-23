//modules/domain_types.go

package modules

type Intent struct {
	Name       string
	Confidence float64
	Domain     string // add this if other code expects it
}

type Task struct {
	Intent Intent
	Action string
}
