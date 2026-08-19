// Package sampling provides transaction dataset sampling strategies for
// approximate mining on large datasets. Random sampling with replacement
// and stratified sampling by transaction size are supported.
package sampling

import (
	"fmt"
	"sort"
)

// Transaction is a basket of items.
type Transaction = []string

// RNG is a simple linear congruential generator for reproducible sampling
// without importing math/rand (which would need a separate source).
type RNG struct {
	state uint64
}

// NewRNG creates a seeded generator.
func NewRNG(seed uint64) *RNG {
	if seed == 0 {
		seed = 0x5DEECE66D
	}
	return &RNG{state: seed}
}

// Next returns the next pseudo-random uint64.
func (r *RNG) Next() uint64 {
	r.state = r.state*6364136223846793005 + 1442695040888963407
	return r.state
}

// Intn returns a value in [0, n).
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.Next() % uint64(n))
}

// RandomSample returns a random sample of size n (with replacement) from
// the transaction set using the given seed.
func RandomSample(txs []Transaction, n int, seed uint64) []Transaction {
	if len(txs) == 0 || n <= 0 {
		return nil
	}
	rng := NewRNG(seed)
	sample := make([]Transaction, n)
	for i := 0; i < n; i++ {
		idx := rng.Intn(len(txs))
		sample[i] = txs[idx]
	}
	return sample
}

// RandomSampleWithoutReplacement returns a sample of size n without
// replacement using Fisher-Yates shuffle on indices.
func RandomSampleWithoutReplacement(txs []Transaction, n int, seed uint64) []Transaction {
	if len(txs) == 0 || n <= 0 {
		return nil
	}
	if n > len(txs) {
		n = len(txs)
	}
	rng := NewRNG(seed)
	indices := make([]int, len(txs))
	for i := range indices {
		indices[i] = i
	}
	// Fisher-Yates
	for i := len(indices) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	sample := make([]Transaction, n)
	for i := 0; i < n; i++ {
		sample[i] = txs[indices[i]]
	}
	return sample
}

// StratifiedSample groups transactions by size (number of items) and
// samples proportionally from each stratum.
func StratifiedSample(txs []Transaction, n int, seed uint64) ([]Transaction, error) {
	if len(txs) == 0 || n <= 0 {
		return nil, nil
	}
	if n > len(txs) {
		return nil, fmt.Errorf("sample size %d exceeds population %d", n, len(txs))
	}

	// Group by transaction size
	strata := make(map[int][]int) // size -> indices
	for i, tx := range txs {
		strata[len(tx)] = append(strata[len(tx)], i)
	}

	// Sort strata keys for determinism
	var sizes []int
	for size := range strata {
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)

	rng := NewRNG(seed)
	var sample []Transaction
	remaining := n

	for i, size := range sizes {
		indices := strata[size]
		// Proportional allocation
		var alloc int
		if i == len(sizes)-1 {
			alloc = remaining
		} else {
			alloc = int(float64(n) * float64(len(indices)) / float64(len(txs)))
			if alloc > remaining {
				alloc = remaining
			}
		}
		if alloc > len(indices) {
			alloc = len(indices)
		}

		// Shuffle and take first alloc
		for k := len(indices) - 1; k > 0; k-- {
			j := rng.Intn(k + 1)
			indices[k], indices[j] = indices[j], indices[k]
		}
		for k := 0; k < alloc; k++ {
			sample = append(sample, txs[indices[k]])
		}
		remaining -= alloc
	}
	return sample, nil
}

// BootstrapSample generates a bootstrap replicate (same size as input,
// with replacement).
func BootstrapSample(txs []Transaction, seed uint64) []Transaction {
	return RandomSample(txs, len(txs), seed)
}
