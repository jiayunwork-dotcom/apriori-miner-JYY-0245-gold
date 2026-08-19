package apriori

import (
	"math"
	"sort"
	"testing"
)

func tx(lines ...string) []Transaction {
	var txs []Transaction
	for _, l := range lines {
		var t Transaction
		for _, f := range splitFields(l) {
			t = append(t, Item(f))
		}
		txs = append(txs, t)
	}
	return txs
}

func splitFields(s string) []string {
	var out []string
	for _, p := range splitWS(s) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func splitWS(s string) []string {
	var parts []string
	cur := ""
	for _, r := range s {
		if r == ' ' || r == '\t' {
			parts = append(parts, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	parts = append(parts, cur)
	return parts
}

func TestApriori(t *testing.T) {
	txs := tx(
		"bread milk",
		"bread diaper beer eggs",
		"milk diaper beer coke",
		"bread milk diaper beer",
		"bread milk diaper coke",
	)
	freq := Apriori(txs, 0.4) // 5 transactions -> threshold ceil(0.4*5)=2
	byKey := map[string]int{}
	for _, f := range freq {
		byKey[key(f.Items)] = f.Count
	}
	check := func(k string, want int) {
		if byKey[k] != want {
			t.Errorf("itemset %q: got %d want %d", k, byKey[k], want)
		}
	}
	// 1-itemsets
	check("bread", 4)
	check("milk", 4)
	check("diaper", 4)
	check("beer", 3)
	// 2-itemsets (keys are sorted items)
	check("bread\x00milk", 3)
	check("bread\x00diaper", 3)
	check("diaper\x00milk", 3)
	// 3-itemsets
	check("bread\x00diaper\x00milk", 2)
	if byKey["beer\x00coke"] != 0 {
		t.Errorf("beer,coke should not be frequent")
	}
}

func TestGenerateRules(t *testing.T) {
	txs := tx(
		"bread milk",
		"bread diaper beer eggs",
		"milk diaper beer coke",
		"bread milk diaper beer",
		"bread milk diaper coke",
	)
	freq := Apriori(txs, 0.4)
	rules := GenerateRules(freq, len(txs), 0.6)
	if len(rules) == 0 {
		t.Fatal("expected at least one rule")
	}
	for _, r := range rules {
		if r.Confidence < 0.6-1e-9 {
			t.Errorf("rule %v=>%v confidence too low: %g", r.Antecedent, r.Consequent, r.Confidence)
		}
		if math.IsNaN(r.Support) || math.IsNaN(r.Confidence) {
			t.Errorf("NaN in rule %+v", r)
		}
	}
}

func TestAprioriEmpty(t *testing.T) {
	result := Apriori(nil, 0.5)
	if len(result) != 0 {
		t.Errorf("empty input should yield no frequent sets, got %d", len(result))
	}
}

func TestAprioriSingleTransaction(t *testing.T) {
	txs := tx("a b c")
	freq := Apriori(txs, 1.0) // threshold=1, everything is frequent
	if len(freq) == 0 {
		t.Fatal("single transaction with support=1.0 should yield frequent sets")
	}
	// should have 1-itemsets: a, b, c; 2-itemsets: ab, ac, bc; 3-itemset: abc
	// total = 3 + 3 + 1 = 7
	if len(freq) != 7 {
		t.Errorf("got %d frequent sets, want 7", len(freq))
	}
}

func TestAprioriHighThreshold(t *testing.T) {
	txs := tx("a b", "c d", "e f")
	freq := Apriori(txs, 0.9) // threshold=3, nothing appears in all 3
	if len(freq) != 0 {
		t.Errorf("high threshold should yield 0, got %d", len(freq))
	}
}

func TestAprioriDuplicateItems(t *testing.T) {
	// duplicate items in a single transaction should be counted once
	txs := tx("a a b b c")
	freq := Apriori(txs, 1.0)
	byKey := map[string]int{}
	for _, f := range freq {
		byKey[key(f.Items)] = f.Count
	}
	if byKey["a"] != 1 {
		t.Errorf("a count = %d, want 1", byKey["a"])
	}
}

func TestAprioriMonotonicity(t *testing.T) {
	// anti-monotone property: if {a,b} is frequent, both {a} and {b} must be
	txs := tx(
		"a b c",
		"a b d",
		"a c d",
		"b c d",
		"a b c d",
	)
	freq := Apriori(txs, 0.4)
	freqKeys := map[string]bool{}
	for _, f := range freq {
		freqKeys[key(f.Items)] = true
	}
	// for every 2-itemset, both 1-subsets must be present
	for _, f := range freq {
		if len(f.Items) == 2 {
			for _, it := range f.Items {
				if !freqKeys[key(Itemset{it})] {
					t.Errorf("anti-monotone violation: %v is frequent but {%s} is not", f.Items, it)
				}
			}
		}
	}
}

func TestGenerateRulesConfidenceBound(t *testing.T) {
	txs := tx("a b", "a b", "a b", "a c", "b c")
	freq := Apriori(txs, 0.4)
	rules := GenerateRules(freq, len(txs), 0.0)
	for _, r := range rules {
		if r.Confidence < 0 || r.Confidence > 1.0+1e-9 {
			t.Errorf("confidence out of [0,1]: %f for %v=>%v", r.Confidence, r.Antecedent, r.Consequent)
		}
	}
}

func TestGenerateRulesSupportConsistency(t *testing.T) {
	txs := tx("x y", "x y z", "x z", "y z", "x y z")
	freq := Apriori(txs, 0.4)
	rules := GenerateRules(freq, len(txs), 0.5)
	for _, r := range rules {
		// support of rule = P(X∪Y) must be <= min(P(X), P(Y))
		if r.Support > 1.0+1e-9 {
			t.Errorf("support > 1: %f", r.Support)
		}
	}
}

func TestGenerateRulesSymmetry(t *testing.T) {
	// for 2-itemset {a,b}: both a=>b and b=>a should be generated (if confidence met)
	txs := tx("a b", "a b", "a b", "a", "b")
	freq := Apriori(txs, 0.4)
	rules := GenerateRules(freq, len(txs), 0.0)
	found := map[string]bool{}
	for _, r := range rules {
		k := joinItems(r.Antecedent) + "=>" + joinItems(r.Consequent)
		found[k] = true
	}
	if !found["a=>b"] || !found["b=>a"] {
		t.Errorf("expected both a=>b and b=>a, got %v", found)
	}
}

func TestItemsetSorting(t *testing.T) {
	items := sorted([]Item{"c", "a", "b"})
	want := Itemset{"a", "b", "c"}
	for i, it := range items {
		if it != want[i] {
			t.Errorf("sorted[%d] = %s, want %s", i, it, want[i])
		}
	}
}

func TestAprioriCountAccuracy(t *testing.T) {
	txs := tx(
		"a b c",
		"a b",
		"a c",
		"b c",
		"a b c",
	)
	freq := Apriori(txs, 0.4)
	byKey := map[string]int{}
	for _, f := range freq {
		byKey[key(f.Items)] = f.Count
	}
	// a appears in tx 0,1,2,4 = 4 times
	if byKey["a"] != 4 {
		t.Errorf("a count = %d, want 4", byKey["a"])
	}
	// a,b appears in tx 0,1,4 = 3 times
	if byKey["a\x00b"] != 3 {
		t.Errorf("a,b count = %d, want 3", byKey["a\x00b"])
	}
	// a,b,c appears in tx 0,4 = 2 times
	if byKey["a\x00b\x00c"] != 2 {
		t.Errorf("a,b,c count = %d, want 2", byKey["a\x00b\x00c"])
	}
}

func joinItems(items Itemset) string {
	parts := make([]string, len(items))
	for i, it := range items {
		parts[i] = string(it)
	}
	sort.Strings(parts)
	return concatParts(parts)
}

func concatParts(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}
