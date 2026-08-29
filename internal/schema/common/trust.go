//internal\schema\common\trust.go

package internal_common

type BootTrust uint8

const (
	TrustInvalid BootTrust = iota
	TrustWeak
	TrustStrong
)
