//modules/domain_types.go

package modules

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil"

type Intent struct {
	Name       string
	Confidence mathutil.Q16
	Domain     string // add this if other code expects it
}

type Task struct {
	Intent Intent
	Action string
}
