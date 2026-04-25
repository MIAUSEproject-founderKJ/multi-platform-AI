//internal/schema/security/permissions.go

package schema_security

import schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"

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

var permissionMap = map[schema_identity.Permission]PermissionMask{
	schema_identity.PermUser:           PermUserMask,
	schema_identity.PermAdmin:          PermAdminMask,
	schema_identity.PermBasicRuntime:   PermBasicRuntimeMask,
	schema_identity.PermConfigEdit:     PermConfigEditMask,
	schema_identity.PermDiagnostics:    PermDiagnosticsMask,
	schema_identity.PermHardwareIO:     PermHardwareIOMask,
	schema_identity.PermSafetyOverride: PermSafetyOverrideMask,
}

//MAP TO MASK
func ToPermissionMask(perms map[schema_identity.Permission]bool) PermissionMask {
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
func FromPermissionMask(mask PermissionMask) map[schema_identity.Permission]bool {
	out := make(map[schema_identity.Permission]bool)

	for p, bit := range permissionMap {
		out[p] = (mask & bit) != 0
	}

	return out
}

func (m PermissionMask) Has(p PermissionMask) bool {
	return m&p != 0
}
