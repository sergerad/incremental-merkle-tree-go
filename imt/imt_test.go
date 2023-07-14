package imt

import (
	"crypto/sha256"
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

	// Instantiate
	imt, err := New(WithHeight(4), WithHash(catFn))
	require.NoError(t, err)

	// Add leaves
	hexRunes := []rune("123456789abcdef")
	leaves := [][]byte{}
	for i := 0; i < imt.Height(); i++ {
		// Get hex characters for leaf
		leaf := "0" + string(hexRunes[i])
		leafBytes, err := hex.DecodeString(leaf)
		require.NoError(t, err)
		t.Log("Adding leaf", leaf)

		// Add leaf
		require.NoError(t, imt.AddLeaf(leafBytes))
		t.Log("Latest root", hex.EncodeToString(imt.RootDigest()))

		// Store leaf for test
		leaves = append(leaves, leafBytes)
	}

	// Verify root
	// Root should be all leaves appended + zeroes
	expectedRoot := []byte{}
	for _, leaf := range leaves {
		expectedRoot = append(expectedRoot, leaf...)
	}
	remainderLeafCount := imt.MaxLeaves() - len(leaves)
	for i := 0; i < remainderLeafCount; i++ {
		b, err := hex.DecodeString("00")
		require.NoError(t, err)
		expectedRoot = append(expectedRoot, b...)
	}
	require.Equal(t, expectedRoot, imt.RootDigest())
}

func TestAddLeaf_WithDefaultValues(t *testing.T) {
	// Instantiate
	imt, err := New()
	require.NoError(t, err)

	// Test default root of a tree of height 32
	expectedRoot := "985e929f70af28d0bdd1a90a808f977f597c7c778c489e98d3bd8910d31ac0f7"
	require.Equal(t, expectedRoot, hex.EncodeToString(imt.RootDigest()))

	// Add leaf
	require.NoError(t, imt.AddLeaf([]byte("test")))
	require.Len(t, imt.RootDigest(), sha256.Size)
}
