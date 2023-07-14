package imt

import (
	"errors"

	"github.com/sergerad/incremental-merkle-tree/hash/sha256"
)

type options struct {
	height     int
	hash       Hash
	digestSize int
}

// Option is a function that sets a non-default
// values for configuration of the Incremental Merkle Tree.
type Option func(*options)

// WithHeight sets a non-default value for the
// height of the tree.
func WithHeight(height int) Option {
	return func(opts *options) {
		opts.height = height
	}
}

// WithHash sets a non-default hashing function
// for the runtime of the tree.
func WithHash(hash Hash) Option {
	return func(opts *options) {
		opts.hash = hash
	}
}

func handleOptions(opts ...Option) (*options, error) {
	// Read options
	o := &options{}
	for _, optsFn := range opts {
		optsFn(o)
	}

	// Default height
	if o.height == 0 {
		o.height = MaxTreeHeight
	}
	// Max height
	if o.height > MaxTreeHeight {
		return nil, ErrTreeHeightTooLarge
	}
	// Default hash
	if o.hash == nil {
		o.hash = sha256.Hash
	}

	// Infer size of digests
	tmpDigest, err := o.hash(make([]byte, 1))
	if err != nil {
		return nil, errors.Join(ErrHashFailed, err)
	}
	o.digestSize = len(tmpDigest)

	return o, nil
}
