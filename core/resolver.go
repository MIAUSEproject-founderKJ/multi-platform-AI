//core/resolver.go
//Capability & Policy Resolver

package core

func HasCapabilities(ctx RuntimeContext, required CapabilitySet) bool {
	return ctx.Capabilities&required == required
}

func HasPermissions(ctx RuntimeContext, required PermissionSet) bool {
	return ctx.Permissions&required == required
}