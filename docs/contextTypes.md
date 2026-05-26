type ExecutionContext interface {
	Platform() internal_environment.PlatformClass
	Capabilities() internal_environment.CapabilitySet
	SecurityTier() user_setting.TrustLevel

	HasPermission(user_setting.PermissionKey) bool
	ServiceType() user_setting.ServiceType
}

type RuntimeContext struct {
	Router     router.Router
	Bus        *bus.MessageBus
	Supervisor *runtime_supervisor.Supervisor

	Modules map[string]runtime_supervisor.Module
}

type RuntimePolicy struct {
	Permissions map[user_setting.PermissionKey]bool
	PermMask    internal_verification.PermissionMask

	AllowNetwork bool
	AllowHotplug bool

	MaxAuthority float64

	AllowedModules map[string]bool
}

type BootContext struct {
	platformClass internal_environment.PlatformClass
	service       user_setting.ServiceType
	entity        internal_environment.EntityKind
	tier          user_setting.TierType
	bootMode      internal_boot.BootMode

	trustLevel  user_setting.TrustLevel
	capabilities internal_environment.CapabilitySet

	vault verification_persistence.VaultStore
	logger *zap.Logger
}