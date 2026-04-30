// internal/schema/boot/boot_mode.go

package internal_boot

type BootMode string

const (
	BootCold BootMode = "cold"
	BootFast BootMode = "fast"
)
