//core/policy/remote_verify.go
package policy
// VerifyRemoteHash checks if the hash provided by a neighbor is authorized.
func (e *TrustEvaluator) VerifyRemoteHash(remoteHash string) bool {
    // 1. Check if the hash matches the local machine's authorized list
    // 2. In a fleet, this ensures all nodes are running the exact same version
    // or a version within the "Trusted Compatibility Range."
    return remoteHash == e.AuthorizedBinaryHash
}