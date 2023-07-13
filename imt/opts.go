package imt

import (
	"github.com/sergerad/incremental-merkle-tree/hash/sha256"
)

type options struct {
	height int
	hash   Hash
}
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

	// Ensure appropriate values
	if o.height == 0 {
		o.height = MaxTreeHeight
	}
	if o.height > MaxTreeHeight {
		return nil, ErrTreeHeightTooLarge
	}
	if o.hash == nil {
		o.hash = sha256.Hash
	}

	return o, nil
}
