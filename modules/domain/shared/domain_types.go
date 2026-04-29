//modules/domain/shared/domain_types.go

package domain_shared

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"

type Intent struct {
	Name       string
	Confidence math_convert.Q16
	Domain     string // add this if other code expects it
}

type Task struct {
	Intent Intent
	Action string
}
