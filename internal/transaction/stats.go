package transaction

import "sort"

// Stats holds summary statistics about a transaction dataset.
type Stats struct {
	Count       int
	TotalItems  int
	UniqueItems int
	MinSize     int
	MaxSize     int
	MeanSize    float64
	MedianSize  float64
}

// ComputeStats calculates summary statistics over a set of transactions.
func ComputeStats(txs []Transaction) Stats {
	if len(txs) == 0 {
		return Stats{}
	}
	s := Stats{Count: len(txs), MinSize: len(txs[0]), MaxSize: len(txs[0])}
	unique := make(map[Item]struct{})
	sizes := make([]int, len(txs))

	for i, tx := range txs {
		size := len(tx)
		sizes[i] = size
		s.TotalItems += size
		if size < s.MinSize {
			s.MinSize = size
		}
		if size > s.MaxSize {
			s.MaxSize = size
		}
		for _, it := range tx {
			unique[it] = struct{}{}
		}
	}

	s.UniqueItems = len(unique)
	s.MeanSize = float64(s.TotalItems) / float64(s.Count)
	s.MedianSize = median(sizes)
	return s
}

// ItemFrequency returns the count of each item across all transactions.
func ItemFrequency(txs []Transaction) map[Item]int {
	freq := make(map[Item]int)
	for _, tx := range txs {
		seen := make(map[Item]bool)
		for _, it := range tx {
			if !seen[it] {
				seen[it] = true
				freq[it]++
			}
		}
	}
	return freq
}

// TopItems returns the n most frequent items sorted by count descending.
func TopItems(txs []Transaction, n int) []ItemCount {
	freq := ItemFrequency(txs)
	items := make([]ItemCount, 0, len(freq))
	for it, c := range freq {
		items = append(items, ItemCount{Item: it, Count: c})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Item < items[j].Item
	})
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

// ItemCount pairs an item with its frequency count.
type ItemCount struct {
	Item  Item
	Count int
}

// RareItems returns items that appear in fewer than minCount transactions.
func RareItems(txs []Transaction, minCount int) []Item {
	freq := ItemFrequency(txs)
	var rare []Item
	for it, c := range freq {
		if c < minCount {
			rare = append(rare, it)
		}
	}
	sort.Strings(rare)
	return rare
}

func median(vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := make([]int, len(vals))
	copy(sorted, vals)
	sort.Ints(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2
	}
	return float64(sorted[n/2])
}
