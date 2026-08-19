package export

import (
	"encoding/json"
	"io"
)

// JSONReport is the structured output for programmatic consumption.
type JSONReport struct {
	Transactions  int          `json:"transactions"`
	MinSupport    float64      `json:"min_support"`
	MinConfidence float64      `json:"min_confidence"`
	FrequentSets  []JSONFreqSet `json:"frequent_sets"`
	Rules         []JSONRule   `json:"rules"`
}

// JSONFreqSet is a frequent itemset in JSON form.
type JSONFreqSet struct {
	Items   []string `json:"items"`
	Count   int      `json:"count"`
	Support float64  `json:"support"`
}

// JSONRule is a rule in JSON form.
type JSONRule struct {
	Antecedent []string `json:"antecedent"`
	Consequent []string `json:"consequent"`
	Support    float64  `json:"support"`
	Confidence float64  `json:"confidence"`
	Lift       float64  `json:"lift"`
	Leverage   float64  `json:"leverage"`
	Conviction float64  `json:"conviction"`
}

// WriteJSONReport writes the complete report as indented JSON.
func WriteJSONReport(w io.Writer, report *JSONReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

// WriteJSONCompact writes the report without indentation.
func WriteJSONCompact(w io.Writer, report *JSONReport) error {
	return json.NewEncoder(w).Encode(report)
}

// WriteJSONLines writes rules as newline-delimited JSON (one per line).
func WriteJSONLines(w io.Writer, rules []JSONRule) error {
	enc := json.NewEncoder(w)
	for _, r := range rules {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}
