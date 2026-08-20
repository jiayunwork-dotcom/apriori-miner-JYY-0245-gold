// Package filter provides rule filtering by metrics thresholds and item constraints.
//
// A Filter is a composable predicate chain: each criterion must pass for a rule
// to be included. Supported criteria:
//   - MinLift / MaxLift
//   - MinLeverage
//   - MinConviction
//   - MustContainAny / MustContainAll (items that must appear in antecedent or consequent)
//   - ExcludeItems (rules containing any excluded item are dropped)
//   - MaxLength (maximum itemset size for antecedent + consequent)
package filter

import "strings"

// Rule represents a rule candidate with its metrics attached.
type Rule struct {
	Antecedent []string
	Consequent []string
	Support    float64
	Confidence float64
	Lift       float64
	Leverage   float64
	Conviction float64
}

// Criteria defines all optional filter parameters.
type Criteria struct {
	MinSupport    *float64
	MaxSupport    *float64
	MinConfidence *float64
	MinLift       *float64
	MaxLift       *float64
	MinLeverage   *float64
	MinConviction *float64
	MustContain   []string // at least one must appear in antecedent or consequent
	MustContainAll []string // all must appear
	ExcludeItems  []string
	MaxLength     *int // max total items in antecedent+consequent
	MinLength     *int // min total items
}

// Apply returns only the rules that pass all criteria.
func Apply(rules []Rule, c Criteria) []Rule {
	var out []Rule
	for _, r := range rules {
		if passes(r, c) {
			out = append(out, r)
		}
	}
	return out
}

// Count returns the number of rules that would pass the criteria.
func Count(rules []Rule, c Criteria) int {
	n := 0
	for _, r := range rules {
		if passes(r, c) {
			n++
		}
	}
	return n
}

// Partition separates rules into (passing, failing) based on criteria.
func Partition(rules []Rule, c Criteria) (pass, fail []Rule) {
	for _, r := range rules {
		if passes(r, c) {
			pass = append(pass, r)
		} else {
			fail = append(fail, r)
		}
	}
	return
}

func passes(r Rule, c Criteria) bool {
	if c.MinSupport != nil && r.Support < *c.MinSupport {
		return false
	}
	if c.MaxSupport != nil && r.Support > *c.MaxSupport {
		return false
	}
	if c.MinConfidence != nil && r.Confidence < *c.MinConfidence {
		return false
	}
	if c.MinLift != nil && r.Lift < *c.MinLift {
		return false
	}
	if c.MaxLift != nil && r.Lift > *c.MaxLift {
		return false
	}
	if c.MinLeverage != nil && r.Leverage < *c.MinLeverage {
		return false
	}
	if c.MinConviction != nil && r.Conviction < *c.MinConviction {
		return false
	}
	if c.MaxLength != nil {
		total := len(r.Antecedent) + len(r.Consequent)
		if total > *c.MaxLength {
			return false
		}
	}
	if c.MinLength != nil {
		total := len(r.Antecedent) + len(r.Consequent)
		if total < *c.MinLength {
			return false
		}
	}
	if len(c.MustContain) > 0 {
		if !containsAny(r, c.MustContain) {
			return false
		}
	}
	if len(c.MustContainAll) > 0 {
		if !containsAllItems(r, c.MustContainAll) {
			return false
		}
	}
	if len(c.ExcludeItems) > 0 {
		if containsExcluded(r, c.ExcludeItems) {
			return false
		}
	}
	return true
}

func containsAny(r Rule, items []string) bool {
	all := allItems(r)
	for _, want := range items {
		for _, got := range all {
			if strings.EqualFold(got, want) {
				return true
			}
		}
	}
	return false
}

func containsAllItems(r Rule, items []string) bool {
	all := allItems(r)
	set := make(map[string]bool)
	for _, it := range all {
		set[strings.ToLower(it)] = true
	}
	for _, want := range items {
		if !set[strings.ToLower(want)] {
			return false
		}
	}
	return true
}

func containsExcluded(r Rule, excluded []string) bool {
	all := allItems(r)
	for _, ex := range excluded {
		for _, it := range all {
			if strings.EqualFold(it, ex) {
				return true
			}
		}
	}
	return false
}

func allItems(r Rule) []string {
	out := make([]string, 0, len(r.Antecedent)+len(r.Consequent))
	out = append(out, r.Antecedent...)
	out = append(out, r.Consequent...)
	return out
}

// TopN returns the top N rules by the given metric, assuming rules are already
// sorted by that metric in descending order.
func TopN(rules []Rule, n int) []Rule {
	if n >= len(rules) {
		return rules
	}
	return rules[:n]
}
