//core/resolver.go
//Capability & Policy Resolver

package core

func hasCapabilities(ctx RuntimeContext, required []Capability) bool {
	for _, cap := range required {
		if !ctx.Platform.Capabilities[cap] {
			return false
		}
	}
	return true
}

func hasPermissions(ctx RuntimeContext, required []string) bool {
	for _, perm := range required {
		if !ctx.Policy.Permissions[perm] {
			return false
		}
	}
	return true
}
