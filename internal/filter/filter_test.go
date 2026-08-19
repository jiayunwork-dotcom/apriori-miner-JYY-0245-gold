package filter

import "testing"

func ptr(f float64) *float64 { return &f }
func intPtr(n int) *int      { return &n }

func sampleRules() []Rule {
	return []Rule{
		{
			Antecedent: []string{"bread"},
			Consequent: []string{"milk"},
			Support:    0.3, Confidence: 0.6, Lift: 1.2, Leverage: 0.05, Conviction: 1.25,
		},
		{
			Antecedent: []string{"beer"},
			Consequent: []string{"diaper"},
			Support:    0.25, Confidence: 0.7, Lift: 1.4, Leverage: 0.08, Conviction: 1.67,
		},
		{
			Antecedent: []string{"bread", "butter"},
			Consequent: []string{"milk"},
			Support:    0.15, Confidence: 0.9, Lift: 1.8, Leverage: 0.12, Conviction: 5.0,
		},
		{
			Antecedent: []string{"eggs"},
			Consequent: []string{"bacon"},
			Support:    0.4, Confidence: 0.5, Lift: 0.9, Leverage: -0.04, Conviction: 0.8,
		},
	}
}

func TestFilterMinLift(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinLift: ptr(1.3)}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (lift >= 1.3)", len(result))
	}
	for _, r := range result {
		if r.Lift < 1.3 {
			t.Errorf("rule %v=>%v has lift %f < 1.3", r.Antecedent, r.Consequent, r.Lift)
		}
	}
}

func TestFilterMaxLift(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MaxLift: ptr(1.3)}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (lift <= 1.3)", len(result))
	}
}

func TestFilterMinSupport(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinSupport: ptr(0.25)}
	result := Apply(rules, c)
	if len(result) != 3 {
		t.Errorf("got %d rules, want 3 (support >= 0.25)", len(result))
	}
}

func TestFilterMinConfidence(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinConfidence: ptr(0.65)}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (confidence >= 0.65)", len(result))
	}
}

func TestFilterMinLeverage(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinLeverage: ptr(0.06)}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (leverage >= 0.06)", len(result))
	}
}

func TestFilterMinConviction(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinConviction: ptr(1.5)}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (conviction >= 1.5)", len(result))
	}
}

func TestFilterMustContain(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MustContain: []string{"bread"}}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (contain bread)", len(result))
	}
}

func TestFilterMustContainAll(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MustContainAll: []string{"bread", "milk"}}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (contain both bread and milk)", len(result))
	}
}

func TestFilterExcludeItems(t *testing.T) {
	rules := sampleRules()
	c := Criteria{ExcludeItems: []string{"beer", "eggs"}}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (exclude beer and eggs)", len(result))
	}
}

func TestFilterMaxLength(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MaxLength: intPtr(2)}
	result := Apply(rules, c)
	if len(result) != 3 {
		t.Errorf("got %d rules, want 3 (max length 2)", len(result))
	}
}

func TestFilterMinLength(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinLength: intPtr(3)}
	result := Apply(rules, c)
	if len(result) != 1 {
		t.Errorf("got %d rules, want 1 (min length 3)", len(result))
	}
}

func TestFilterCombined(t *testing.T) {
	rules := sampleRules()
	c := Criteria{
		MinLift:       ptr(1.0),
		MinConfidence: ptr(0.6),
		MustContain:   []string{"milk"},
	}
	result := Apply(rules, c)
	if len(result) != 2 {
		t.Errorf("got %d rules, want 2 (combined)", len(result))
	}
}

func TestFilterNoCriteria(t *testing.T) {
	rules := sampleRules()
	c := Criteria{}
	result := Apply(rules, c)
	if len(result) != len(rules) {
		t.Errorf("no criteria should return all %d rules, got %d", len(rules), len(result))
	}
}

func TestCount(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinLift: ptr(1.3)}
	n := Count(rules, c)
	if n != 2 {
		t.Errorf("Count = %d, want 2", n)
	}
}

func TestPartition(t *testing.T) {
	rules := sampleRules()
	c := Criteria{MinLift: ptr(1.0)}
	pass, fail := Partition(rules, c)
	if len(pass) != 3 {
		t.Errorf("pass = %d, want 3", len(pass))
	}
	if len(fail) != 1 {
		t.Errorf("fail = %d, want 1", len(fail))
	}
}

func TestTopN(t *testing.T) {
	rules := sampleRules()
	top := TopN(rules, 2)
	if len(top) != 2 {
		t.Errorf("TopN(2) = %d rules, want 2", len(top))
	}
}

func TestTopNLargerThanSlice(t *testing.T) {
	rules := sampleRules()
	top := TopN(rules, 100)
	if len(top) != len(rules) {
		t.Errorf("TopN(100) = %d rules, want %d", len(top), len(rules))
	}
}

func TestFilterCaseInsensitive(t *testing.T) {
	rules := []Rule{
		{Antecedent: []string{"Bread"}, Consequent: []string{"MILK"}, Lift: 1.5},
	}
	c := Criteria{MustContain: []string{"bread"}}
	result := Apply(rules, c)
	if len(result) != 1 {
		t.Error("case-insensitive matching should find Bread")
	}
}

func TestFilterExcludeCaseInsensitive(t *testing.T) {
	rules := []Rule{
		{Antecedent: []string{"Beer"}, Consequent: []string{"chips"}, Lift: 1.1},
		{Antecedent: []string{"bread"}, Consequent: []string{"milk"}, Lift: 1.2},
	}
	c := Criteria{ExcludeItems: []string{"BEER"}}
	result := Apply(rules, c)
	if len(result) != 1 {
		t.Errorf("got %d rules, want 1 (exclude BEER case-insensitive)", len(result))
	}
}
