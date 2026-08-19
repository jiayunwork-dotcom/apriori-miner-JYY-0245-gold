package rank

import "sort"

// DiverseTopK selects top-K rules while maximising diversity: no two
// selected rules should share more than maxOverlap fraction of their
// items. This prevents the result set from being dominated by trivially
// similar rules (e.g., {A,B}=>{C} and {A,B,D}=>{C}).
func DiverseTopK(rules []Rule, k int, maxOverlap float64) []Rule {
	if k >= len(rules) {
		return rules
	}
	// Work on a sorted copy (by confidence desc)
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Confidence > sorted[j].Confidence
	})

	var selected []Rule
	for _, r := range sorted {
		if len(selected) >= k {
			break
		}
		tooSimilar := false
		rItems := allItemsSet(r)
		for _, s := range selected {
			sItems := allItemsSet(s)
			overlap := jaccardSimilarity(rItems, sItems)
			if overlap > maxOverlap {
				tooSimilar = true
				break
			}
		}
		if !tooSimilar {
			selected = append(selected, r)
		}
	}
	return selected
}

func allItemsSet(r Rule) map[string]bool {
	s := make(map[string]bool)
	for _, it := range r.Antecedent {
		s[it] = true
	}
	for _, it := range r.Consequent {
		s[it] = true
	}
	return s
}

func jaccardSimilarity(a, b map[string]bool) float64 {
	inter := 0
	for k := range a {
		if b[k] {
			inter++
		}
	}
	union := len(a) + len(b) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

// CoverageScore computes what fraction of unique items in the dataset
// are mentioned in the given rule set. Higher coverage means the rules
// span more of the item space.
func CoverageScore(rules []Rule, totalItems int) float64 {
	covered := make(map[string]bool)
	for _, r := range rules {
		for _, it := range r.Antecedent {
			covered[it] = true
		}
		for _, it := range r.Consequent {
			covered[it] = true
		}
	}
	if totalItems == 0 {
		return 0
	}
	return float64(len(covered)) / float64(totalItems)
}

// Novelty scores a rule based on how different it is from the rest of
// the rule set. A highly novel rule shares few items with other rules.
func Novelty(r Rule, others []Rule) float64 {
	if len(others) == 0 {
		return 1
	}
	rItems := allItemsSet(r)
	sumSim := 0.0
	for _, o := range others {
		oItems := allItemsSet(o)
		sumSim += jaccardSimilarity(rItems, oItems)
	}
	avgSim := sumSim / float64(len(others))
	return 1 - avgSim
}
