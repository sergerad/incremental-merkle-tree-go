package sha256

import (
	"crypto/sha256"
)

// Hash returns the SHA256 digest of the given data
func Hash(data ...[]byte) ([]byte, error) {
	hasher := sha256.New()
	for i := 0; i < len(data); i++ {
		if _, err := hasher.Write(data[i]); err != nil {
			return nil, err
		}
	}
	return hasher.Sum(nil), nil
}
