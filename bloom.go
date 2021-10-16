package dht

import (
	"github.com/bits-and-blooms/bloom"
)

func newBloomFilterForTraversal() *bloom.BloomFilter {
	return bloom.NewWithEstimates(10000, 0.5)
}
