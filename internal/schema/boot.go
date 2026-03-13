//internal/schema/boot.go

package schema

type BootMode string

const (
	BootCold BootMode = "cold"
	BootFast BootMode = "fast"
)
