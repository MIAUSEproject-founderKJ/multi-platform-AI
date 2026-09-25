// core/security/persistence/golden_hash_store.go

package verification_persistence

func (v *IsolatedVault) SealGoldenHash(machine string, hash []byte) error {
	return v.writeEncrypted("golden-"+machine, hash)
}

func (v *IsolatedVault) LoadGoldenHash(machine string) (string, error) {
	data, err := v.readDecrypted("golden-" + machine)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
