package filter

import "sort"

// SortByConfidence sorts rules by confidence descending.
func SortByConfidence(rules []Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Confidence > rules[j].Confidence
	})
}

// SortByLift sorts rules by lift descending.
func SortByLift(rules []Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Lift > rules[j].Lift
	})
}

// SortBySupport sorts rules by support descending.
func SortBySupport(rules []Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Support > rules[j].Support
	})
}

// SortByLeverage sorts rules by leverage descending.
func SortByLeverage(rules []Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Leverage > rules[j].Leverage
	})
}

// UniqueRules removes duplicate rules (same antecedent and consequent).
func UniqueRules(rules []Rule) []Rule {
	seen := make(map[string]bool)
	var out []Rule
	for _, r := range rules {
		k := ruleKey(r)
		if !seen[k] {
			seen[k] = true
			out = append(out, r)
		}
	}
	return out
}

func ruleKey(r Rule) string {
	ante := make([]string, len(r.Antecedent))
	copy(ante, r.Antecedent)
	sort.Strings(ante)
	cons := make([]string, len(r.Consequent))
	copy(cons, r.Consequent)
	sort.Strings(cons)
	key := ""
	for _, a := range ante {
		key += a + "|"
	}
	key += "=>"
	for _, c := range cons {
		key += c + "|"
	}
	return key
}

// GroupByConsequent groups rules by their consequent items.
func GroupByConsequent(rules []Rule) map[string][]Rule {
	groups := make(map[string][]Rule)
	for _, r := range rules {
		cons := make([]string, len(r.Consequent))
		copy(cons, r.Consequent)
		sort.Strings(cons)
		key := ""
		for _, c := range cons {
			key += c + ","
		}
		groups[key] = append(groups[key], r)
	}
	return groups
}

// GroupByAntecedent groups rules by their antecedent items.
func GroupByAntecedent(rules []Rule) map[string][]Rule {
	groups := make(map[string][]Rule)
	for _, r := range rules {
		ante := make([]string, len(r.Antecedent))
		copy(ante, r.Antecedent)
		sort.Strings(ante)
		key := ""
		for _, a := range ante {
			key += a + ","
		}
		groups[key] = append(groups[key], r)
	}
	return groups
}
