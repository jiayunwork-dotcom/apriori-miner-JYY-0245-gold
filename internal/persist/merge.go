package persist

import (
	"fmt"
	"sort"
	"strings"
)

// Merge combines two mining results into one. Frequent sets are unioned by
// itemset key; counts are summed. Rules are de-duplicated by antecedent/consequent
// key; metrics from the result with higher support win.
//
// This supports incremental mining: mine batches separately, then merge results
// before persisting a combined view.
func Merge(a, b *MiningResult) *MiningResult {
	merged := &MiningResult{
		Version:           currentVersion,
		TotalTransactions: a.TotalTransactions + b.TotalTransactions,
		MinSupport:        minFloat(a.MinSupport, b.MinSupport),
		MinConfidence:     minFloat(a.MinConfidence, b.MinConfidence),
	}

	// merge frequent sets
	fsMap := make(map[string]*FrequentSetRecord)
	for i := range a.FrequentSets {
		fs := a.FrequentSets[i]
		k := itemsetKey(fs.Items)
		fsMap[k] = &FrequentSetRecord{
			Items:   fs.Items,
			Count:   fs.Count,
			Support: fs.Support,
		}
	}
	for i := range b.FrequentSets {
		fs := b.FrequentSets[i]
		k := itemsetKey(fs.Items)
		if existing, ok := fsMap[k]; ok {
			existing.Count += fs.Count
			// recalculate support based on merged total
			existing.Support = float64(existing.Count) / float64(merged.TotalTransactions)
		} else {
			fsMap[k] = &FrequentSetRecord{
				Items:   fs.Items,
				Count:   fs.Count,
				Support: float64(fs.Count) / float64(merged.TotalTransactions),
			}
		}
	}

	merged.FrequentSets = make([]FrequentSetRecord, 0, len(fsMap))
	for _, fs := range fsMap {
		merged.FrequentSets = append(merged.FrequentSets, *fs)
	}
	sort.Slice(merged.FrequentSets, func(i, j int) bool {
		return merged.FrequentSets[i].Count > merged.FrequentSets[j].Count
	})

	// merge rules
	ruleMap := make(map[string]*RuleRecord)
	for i := range a.Rules {
		r := a.Rules[i]
		k := ruleKey(r)
		ruleMap[k] = &r
	}
	for i := range b.Rules {
		r := b.Rules[i]
		k := ruleKey(r)
		if existing, ok := ruleMap[k]; ok {
			// keep the one with higher support
			if r.Support > existing.Support {
				ruleMap[k] = &r
			}
		} else {
			ruleMap[k] = &r
		}
	}

	merged.Rules = make([]RuleRecord, 0, len(ruleMap))
	for _, r := range ruleMap {
		merged.Rules = append(merged.Rules, *r)
	}
	sort.Slice(merged.Rules, func(i, j int) bool {
		return merged.Rules[i].Confidence > merged.Rules[j].Confidence
	})

	return merged
}

// Summary returns a human-readable summary of the mining result.
func Summary(r *MiningResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Mining Result Summary\n")
	fmt.Fprintf(&b, "  Transactions: %d\n", r.TotalTransactions)
	fmt.Fprintf(&b, "  Min Support: %.4f\n", r.MinSupport)
	fmt.Fprintf(&b, "  Min Confidence: %.4f\n", r.MinConfidence)
	fmt.Fprintf(&b, "  Frequent Sets: %d\n", len(r.FrequentSets))
	fmt.Fprintf(&b, "  Rules: %d\n", len(r.Rules))

	if len(r.FrequentSets) > 0 {
		maxSize := 0
		for _, fs := range r.FrequentSets {
			if len(fs.Items) > maxSize {
				maxSize = len(fs.Items)
			}
		}
		fmt.Fprintf(&b, "  Max Itemset Size: %d\n", maxSize)
	}

	if len(r.Rules) > 0 {
		maxConf := 0.0
		for _, rule := range r.Rules {
			if rule.Confidence > maxConf {
				maxConf = rule.Confidence
			}
		}
		fmt.Fprintf(&b, "  Max Confidence: %.4f\n", maxConf)
	}

	return b.String()
}

// Validate checks basic integrity constraints on a mining result.
func Validate(r *MiningResult) error {
	if r.TotalTransactions < 0 {
		return fmt.Errorf("persist: negative transaction count: %d", r.TotalTransactions)
	}
	if r.MinSupport < 0 || r.MinSupport > 1 {
		return fmt.Errorf("persist: min_support out of [0,1]: %f", r.MinSupport)
	}
	if r.MinConfidence < 0 || r.MinConfidence > 1 {
		return fmt.Errorf("persist: min_confidence out of [0,1]: %f", r.MinConfidence)
	}
	for i, fs := range r.FrequentSets {
		if fs.Count < 0 {
			return fmt.Errorf("persist: frequent_set[%d] has negative count", i)
		}
		if len(fs.Items) == 0 {
			return fmt.Errorf("persist: frequent_set[%d] has no items", i)
		}
	}
	for i, rule := range r.Rules {
		if len(rule.Antecedent) == 0 {
			return fmt.Errorf("persist: rule[%d] has empty antecedent", i)
		}
		if len(rule.Consequent) == 0 {
			return fmt.Errorf("persist: rule[%d] has empty consequent", i)
		}
		if rule.Confidence < 0 || rule.Confidence > 1+1e-9 {
			return fmt.Errorf("persist: rule[%d] confidence out of [0,1]: %f", i, rule.Confidence)
		}
	}
	return nil
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
