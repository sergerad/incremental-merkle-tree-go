package imt

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddLeaf_WithAppendHasher(t *testing.T) {
	// Append inputs for testing
	catFn := func(data ...[]byte) ([]byte, error) {
		if len(data) == 1 {
			return data[0], nil
		}
		return append(data[0], data[1]...), nil
	}

	height := 4
	imt, err := New(height, 2, catFn)
	require.NoError(t, err)

	hexRunes := []rune("123456789abcdef")
	for i := 0; i < height; i++ {
		// Get hex characters for leaf
		leaf := "0" + string(hexRunes[i])
		leafBytes, err := hex.DecodeString(leaf)
		require.NoError(t, err)
		t.Log("Added leaf", leaf)

		// Add leaf
		require.NoError(t, imt.AddLeaf(leafBytes))
		t.Log("Latest root", hex.EncodeToString(imt.RootDigest()))
	}
}
