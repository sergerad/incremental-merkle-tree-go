package imt

import (
	"errors"
	"fmt"
	"math"
)

const (
	// MaxTreeHeight is the maximum height of the tree.
	// 32 levels is sufficient to support 2^32 leaves.
	MaxTreeHeight = 32
)

var (
	ErrTreeHeightTooLarge = errors.New(fmt.Sprint("tree height is too large, max is ", MaxTreeHeight))
	ErrTreeIsFull         = errors.New("tree is at maximum leaf capacity")
	ErrHashFailed         = errors.New("failed to perform hash function")
)

// Hash represents a function that creates a digest
// from some provided data.
type Hash func(data ...[]byte) ([]byte, error)

// IncrementalMerkleTree is an incremental Merkle tree.
// It can be used to continuously calculate a Merkle
// root hash in polynomial time as the leaves are
// added to the tree.
type IncrementalMerkleTree struct {
	options

	leftDigestsPerLevel [][]byte
	zeroDigestsPerLevel [][]byte
	rootDigest          []byte

	maxLeaves int

	nextLeafIndex int
}

// New instantiates an Incremental Merkle Tree.
// The height of the tree determines the maximum number of leaves
// that can be added to the tree (2^height).
func New(opts ...Option) (*IncrementalMerkleTree, error) {
	// Process options
	o, err := handleOptions(opts...)
	if err != nil {
		return nil, err
	}

	// Infer size of digests
	tmpDigest, err := o.hash(make([]byte, 1))
	if err != nil {
		return nil, errors.Join(ErrHashFailed, err)
	}

	// Create all zero digests
	zeroDigests := make([][]byte, o.height)
	zeroDigests[0] = make([]byte, len(tmpDigest))
	for i := 1; i < o.height; i++ {
		digest, err := o.hash(zeroDigests[i-1], zeroDigests[i-1])
		if err != nil {
			return nil, errors.Join(ErrHashFailed, err)
		}
		zeroDigests[i] = digest
	}

	return &IncrementalMerkleTree{
		options:             *o,
		leftDigestsPerLevel: make([][]byte, o.height),
		zeroDigestsPerLevel: zeroDigests,
		maxLeaves:           int(math.Pow(2, float64(o.height))),
	}, nil
}

func (imt *IncrementalMerkleTree) AddLeaf(leaf []byte) error {
	// Cannot add more leaves than the height of the tree allows for
	if imt.nextLeafIndex >= imt.maxLeaves {
		return ErrTreeIsFull
	}

	// Start the index at the expected next leaf index.
	// We will use this index to traverse the tree nodes
	// upwards to the root.
	leftRightIndex := imt.nextLeafIndex
	latestDigest, err := imt.hash(leaf)
	if err != nil {
		return errors.Join(ErrHashFailed, err)
	}

	// Iterate through the levels of the tree,
	// starting from the bottom.
	for level := 0; level < imt.height; level++ {
		// We want to hash two nodes together
		var left, right []byte
		// Determine which nodes to hash together based on
		// our current position in the tree.
		// If the index is even, we are on a left node.
		if leftRightIndex%2 == 0 {
			// The left digest is the digest from the
			// the previous level (or the leaf itself).
			left = latestDigest
			// Right is always the zero digest
			right = imt.zeroDigestsPerLevel[level]
			// For every new leaf, we update the list of
			// digests for each level of the tree.
			imt.leftDigestsPerLevel[level] = left
		} else {
			// Left was calculated in previous executions
			// of leaf addition.
			left = imt.leftDigestsPerLevel[level]
			// Right is the digest from the last level or
			// the leaf itself.
			right = latestDigest
		}
		// Append left and right and hash them together
		latestDigest, err = imt.hash(left, right)
		if err != nil {
			return errors.Join(ErrHashFailed, err)
		}
		// Divide the index by 2 to traverse up the tree.
		// E.G. (0,1)->L, (2,3)->R, and so on.
		leftRightIndex /= 2

	}

	// Store the new root digest
	imt.rootDigest = latestDigest

	// Iterate the index so that we can tell
	// whether the next leaf is a left or right.
	imt.nextLeafIndex++

	return nil
}

// RootDigest returns the current root of the tree.
// The root digest changes every time a leaf is added to the tree.
func (imt *IncrementalMerkleTree) RootDigest() []byte {
	root := make([]byte, len(imt.rootDigest))
	copy(root, imt.rootDigest)
	return root
}

// Height returns the number of levels to the tree.
// Height is set on creation and cannot be changed.
func (imt *IncrementalMerkleTree) Height() int {
	return imt.height
}

// MaxLeaves returns the maximum number of leaves that can be
// added to the tree. This is determined by the height of the tree.
func (imt *IncrementalMerkleTree) MaxLeaves() int {
	return imt.maxLeaves
}
