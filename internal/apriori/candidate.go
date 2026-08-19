package apriori

import "sort"

// Candidate is a potential frequent itemset being evaluated.
type Candidate struct {
	Items Itemset
	Count int
}

// GenerateCandidates produces (k)-itemset candidates from (k-1) frequent
// itemsets using the join step of the Apriori algorithm.
func GenerateCandidates(prev []Itemset, k int) []Itemset {
	return join(prev, k)
}

// PruneCandidates removes candidates whose (k-1)-subsets are not all in
// the frequent set. Returns only valid candidates.
func PruneCandidates(candidates []Itemset, frequentKeys map[string]bool) []Itemset {
	var pruned []Itemset
	for _, c := range candidates {
		if hasFrequentSubsets(c, frequentKeys) {
			pruned = append(pruned, c)
		}
	}
	return pruned
}

// CountSupport counts how many transactions in txs contain each candidate.
func CountSupport(txs []Transaction, candidates []Itemset) map[string]int {
	counts := make(map[string]int)
	for _, tx := range txs {
		txset := make(map[Item]bool)
		for _, it := range tx {
			txset[it] = true
		}
		for _, c := range candidates {
			if containsAll(txset, c) {
				counts[key(c)]++
			}
		}
	}
	return counts
}

// FilterByThreshold keeps only candidates meeting the minimum count.
func FilterByThreshold(candidates []Itemset, counts map[string]int, threshold int) []FrequentSet {
	var result []FrequentSet
	for _, c := range candidates {
		if counts[key(c)] >= threshold {
			result = append(result, FrequentSet{Items: c, Count: counts[key(c)]})
		}
	}
	return result
}

// SortFrequentSets sorts frequent sets by count descending, then by key
// ascending for stability.
func SortFrequentSets(sets []FrequentSet) {
	sort.Slice(sets, func(i, j int) bool {
		if sets[i].Count != sets[j].Count {
			return sets[i].Count > sets[j].Count
		}
		return key(sets[i].Items) < key(sets[j].Items)
	})
}

// MaxItemsetSize returns the size of the largest frequent itemset.
func MaxItemsetSize(sets []FrequentSet) int {
	max := 0
	for _, fs := range sets {
		if len(fs.Items) > max {
			max = len(fs.Items)
		}
	}
	return max
}

// FrequentOfSize returns only itemsets of exactly the given size.
func FrequentOfSize(sets []FrequentSet, size int) []FrequentSet {
	var out []FrequentSet
	for _, fs := range sets {
		if len(fs.Items) == size {
			out = append(out, fs)
		}
	}
	return out
}
