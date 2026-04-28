//internal/schema/verification/permissions.go

package internal_verification

import user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"

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

var permissionMap = map[user_setting.Permission]PermissionMask{
	user_setting.PermUser:           PermUserMask,
	user_setting.PermAdmin:          PermAdminMask,
	user_setting.PermBasicRuntime:   PermBasicRuntimeMask,
	user_setting.PermConfigEdit:     PermConfigEditMask,
	user_setting.PermDiagnostics:    PermDiagnosticsMask,
	user_setting.PermHardwareIO:     PermHardwareIOMask,
	user_setting.PermSafetyOverride: PermSafetyOverrideMask,
}

//MAP TO MASK
func ToPermissionMask(perms map[user_setting.Permission]bool) PermissionMask {
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
func FromPermissionMask(mask PermissionMask) map[user_setting.Permission]bool {
	out := make(map[user_setting.Permission]bool)

	for p, bit := range permissionMap {
		out[p] = (mask & bit) != 0
	}

	return out
}

func (m PermissionMask) Has(p PermissionMask) bool {
	return m&p != 0
}
