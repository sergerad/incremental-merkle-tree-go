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

type hashFn = func(data ...[]byte) ([]byte, error)

// IncrementalMerkleTree is an incremental Merkle tree.
// It can be used to continuously calculate a Merkle
// root hash in polynomial time as the leaves are
// added to the tree.
type IncrementalMerkleTree struct {
	leftDigestsPerLevel [][]byte
	zeroDigestsPerLevel [][]byte
	rootDigest          []byte

	hasher    hashFn
	height    int
	maxLeaves int

	nextLeafIndex int
}

func New(height, digestLen int, hasher hashFn) (*IncrementalMerkleTree, error) {
	if height > MaxTreeHeight {
		return nil, ErrTreeHeightTooLarge
	}

	zeroDigests := make([][]byte, height)
	zeroDigests[0] = make([]byte, digestLen)
	for i := 1; i < height; i++ {
		digest, err := hasher(zeroDigests[i-1], zeroDigests[i-1])
		if err != nil {
			return nil, errors.Join(ErrHashFailed, err)
		}
		zeroDigests[i] = digest
	}

	// First root is result of merkling all zeros
	return &IncrementalMerkleTree{
		leftDigestsPerLevel: make([][]byte, height),
		zeroDigestsPerLevel: zeroDigests,
		hasher:              hasher,
		height:              height,
		maxLeaves:           int(math.Pow(2, float64(height))),
	}, nil
}

func (imt *IncrementalMerkleTree) AddLeaf(leaf []byte) error {
	if imt.nextLeafIndex >= imt.maxLeaves {
		return ErrTreeIsFull
	}

	// Start the index at the expected next leaf index.
	// We will use this index to traverse the tree nodes
	// upwards to the root.
	leftRightIndex := imt.nextLeafIndex
	latestDigest, err := imt.hasher(leaf)
	if err != nil {
		return errors.Join(ErrHashFailed, err)
	}
	//println("xxxxxxxxx")
	//println("initial: ", hex.EncodeToString(latestDigest))
	//println("nextindex: ", leftRightIndex)

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
		latestDigest, err = imt.hasher(left, right)
		if err != nil {
			return errors.Join(ErrHashFailed, err)
		}
		// Divide the index by 2 to traverse up the tree.
		// E.G. (0,1)->L, (2,3)->R, and so on.
		leftRightIndex /= 2

		//println("current: ", hex.EncodeToString(latestDigest))
	}

	// Store the new root digest
	imt.rootDigest = latestDigest
	// Iterate the index so that we can tell
	// whether the next leaf is a left or right.
	imt.nextLeafIndex++
	//println("root", hex.EncodeToString(imt.rootDigest))

	return nil
}

func (imt *IncrementalMerkleTree) RootDigest() []byte {
	root := make([]byte, len(imt.rootDigest))
	copy(root, imt.rootDigest)
	return root
}