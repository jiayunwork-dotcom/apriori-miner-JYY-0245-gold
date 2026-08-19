package persist

import (
	"testing"
)

func TestMergeFrequentSets(t *testing.T) {
	a := &MiningResult{
		TotalTransactions: 100,
		MinSupport:        0.2,
		MinConfidence:     0.6,
		FrequentSets: []FrequentSetRecord{
			{Items: []string{"bread"}, Count: 60, Support: 0.6},
			{Items: []string{"milk"}, Count: 40, Support: 0.4},
		},
	}
	b := &MiningResult{
		TotalTransactions: 50,
		MinSupport:        0.3,
		MinConfidence:     0.5,
		FrequentSets: []FrequentSetRecord{
			{Items: []string{"bread"}, Count: 30, Support: 0.6},
			{Items: []string{"butter"}, Count: 20, Support: 0.4},
		},
	}

	merged := Merge(a, b)

	if merged.TotalTransactions != 150 {
		t.Errorf("TotalTransactions = %d, want 150", merged.TotalTransactions)
	}
	if merged.MinSupport != 0.2 {
		t.Errorf("MinSupport = %f, want 0.2", merged.MinSupport)
	}

	fsMap := make(map[string]int)
	for _, fs := range merged.FrequentSets {
		fsMap[itemsetKey(fs.Items)] = fs.Count
	}

	if fsMap["bread"] != 90 {
		t.Errorf("bread count = %d, want 90", fsMap["bread"])
	}
	if fsMap["milk"] != 40 {
		t.Errorf("milk count = %d, want 40", fsMap["milk"])
	}
	if fsMap["butter"] != 20 {
		t.Errorf("butter count = %d, want 20", fsMap["butter"])
	}
}

func TestMergeRules(t *testing.T) {
	a := &MiningResult{
		TotalTransactions: 100,
		MinSupport:        0.2,
		MinConfidence:     0.5,
		Rules: []RuleRecord{
			{Antecedent: []string{"bread"}, Consequent: []string{"milk"}, Support: 0.3, Confidence: 0.5},
		},
	}
	b := &MiningResult{
		TotalTransactions: 100,
		MinSupport:        0.2,
		MinConfidence:     0.5,
		Rules: []RuleRecord{
			{Antecedent: []string{"bread"}, Consequent: []string{"milk"}, Support: 0.4, Confidence: 0.7},
			{Antecedent: []string{"beer"}, Consequent: []string{"chips"}, Support: 0.2, Confidence: 0.6},
		},
	}

	merged := Merge(a, b)
	if len(merged.Rules) != 2 {
		t.Fatalf("rules count = %d, want 2", len(merged.Rules))
	}

	// bread=>milk should keep higher support version
	for _, r := range merged.Rules {
		if r.Antecedent[0] == "bread" && r.Consequent[0] == "milk" {
			if r.Support != 0.4 {
				t.Errorf("bread=>milk support = %f, want 0.4", r.Support)
			}
		}
	}
}

func TestMergeEmpty(t *testing.T) {
	a := &MiningResult{TotalTransactions: 10, MinSupport: 0.1, MinConfidence: 0.5}
	b := &MiningResult{TotalTransactions: 0, MinSupport: 0.2, MinConfidence: 0.6}
	merged := Merge(a, b)
	if merged.TotalTransactions != 10 {
		t.Errorf("TotalTransactions = %d, want 10", merged.TotalTransactions)
	}
}

func TestSummary(t *testing.T) {
	r := sampleResult()
	s := Summary(r)
	if len(s) == 0 {
		t.Error("Summary returned empty string")
	}
}

func TestValidateValid(t *testing.T) {
	r := sampleResult()
	if err := Validate(r); err != nil {
		t.Errorf("valid result failed validation: %v", err)
	}
}

func TestValidateNegativeCount(t *testing.T) {
	r := sampleResult()
	r.FrequentSets[0].Count = -1
	if err := Validate(r); err == nil {
		t.Error("expected error for negative count")
	}
}

func TestValidateEmptyAntecedent(t *testing.T) {
	r := sampleResult()
	r.Rules[0].Antecedent = nil
	if err := Validate(r); err == nil {
		t.Error("expected error for empty antecedent")
	}
}

func TestValidateBadConfidence(t *testing.T) {
	r := sampleResult()
	r.Rules[0].Confidence = 1.5
	if err := Validate(r); err == nil {
		t.Error("expected error for confidence > 1")
	}
}

func TestValidateBadMinSupport(t *testing.T) {
	r := &MiningResult{MinSupport: -0.1}
	if err := Validate(r); err == nil {
		t.Error("expected error for negative min_support")
	}
}
