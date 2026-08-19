package export

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteRulesCSV(t *testing.T) {
	rules := []RuleRow{
		{Antecedent: "bread", Consequent: "milk", Support: 0.5, Confidence: 0.8, Lift: 1.6},
	}
	var buf bytes.Buffer
	err := WriteRulesCSV(&buf, rules)
	if err != nil {
		t.Fatalf("WriteRulesCSV: %v", err)
	}
	if !strings.Contains(buf.String(), "antecedent") {
		t.Error("missing header")
	}
	if !strings.Contains(buf.String(), "bread") {
		t.Error("missing data")
	}
}

func TestWriteJSONReport(t *testing.T) {
	report := &JSONReport{
		Transactions:  100,
		MinSupport:    0.2,
		MinConfidence: 0.6,
		Rules: []JSONRule{
			{Antecedent: []string{"a"}, Consequent: []string{"b"}, Support: 0.3, Confidence: 0.7},
		},
	}
	var buf bytes.Buffer
	err := WriteJSONReport(&buf, report)
	if err != nil {
		t.Fatalf("WriteJSONReport: %v", err)
	}
	if !strings.Contains(buf.String(), "\"transactions\"") {
		t.Error("missing transactions field")
	}
}
