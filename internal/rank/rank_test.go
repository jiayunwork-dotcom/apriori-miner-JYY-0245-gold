package rank

import "testing"

func TestSortByLift(t *testing.T) {
	rules := []Rule{
		{Antecedent: []string{"a"}, Consequent: []string{"b"}, Lift: 2.0},
		{Antecedent: []string{"c"}, Consequent: []string{"d"}, Lift: 5.0},
		{Antecedent: []string{"e"}, Consequent: []string{"f"}, Lift: 1.0},
	}
	SortBy(rules, ByLift)
	if rules[0].Lift != 5.0 {
		t.Errorf("first rule lift = %f, want 5.0", rules[0].Lift)
	}
	if rules[2].Lift != 1.0 {
		t.Errorf("last rule lift = %f, want 1.0", rules[2].Lift)
	}
}

func TestDiverseTopK(t *testing.T) {
	rules := []Rule{
		{Antecedent: []string{"a", "b"}, Consequent: []string{"c"}, Confidence: 0.9},
		{Antecedent: []string{"a", "b", "d"}, Consequent: []string{"c"}, Confidence: 0.85},
		{Antecedent: []string{"x"}, Consequent: []string{"y"}, Confidence: 0.8},
	}
	// With low overlap threshold, similar rules should be excluded
	selected := DiverseTopK(rules, 2, 0.3)
	if len(selected) != 2 {
		t.Fatalf("got %d rules, want 2", len(selected))
	}
	// First should be highest confidence
	if selected[0].Confidence != 0.9 {
		t.Errorf("first conf = %f, want 0.9", selected[0].Confidence)
	}
}
