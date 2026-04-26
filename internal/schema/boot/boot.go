//internal/schema/boot/boot.go
package schema_boot

type BootMode string

const (
	BootCold BootMode = "cold"
	BootFast BootMode = "fast"
)
