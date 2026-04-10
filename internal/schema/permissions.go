//internal/schema/permissions.go

package schema

type PermissionMask uint64

const (
	PermUserMask PermissionMask = 1 << iota
	PermAdminMask
	PermBasicRuntimeMask
	PermConfigEditMask
	PermDiagnosticsMask
	PermHardwareIOMask
	PermSafetyOverrideMask
)

var permissionMap = map[Permission]PermissionMask{
	PermUser:           PermUserMask,
	PermAdmin:          PermAdminMask,
	PermBasicRuntime:   PermBasicRuntimeMask,
	PermConfigEdit:     PermConfigEditMask,
	PermDiagnostics:    PermDiagnosticsMask,
	PermHardwareIO:     PermHardwareIOMask,
	PermSafetyOverride: PermSafetyOverrideMask,
}

//MAP TO MASK
func ToPermissionMask(perms map[Permission]bool) PermissionMask {
	var mask PermissionMask

	for p, enabled := range perms {
		if !enabled {
			continue
		}
		if bit, ok := permissionMap[p]; ok {
			mask |= bit
		}
	}

	return mask
}

//MASK TO MAP
func FromPermissionMask(mask PermissionMask) map[Permission]bool {
	out := make(map[Permission]bool)

	for p, bit := range permissionMap {
		out[p] = (mask & bit) != 0
	}

	return out
}

func (m PermissionMask) Has(p PermissionMask) bool {
	return m&p != 0
}