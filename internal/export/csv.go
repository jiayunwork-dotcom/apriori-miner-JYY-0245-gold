// Package export serialises mining results into interchange formats: CSV
// for spreadsheet consumption and JSON for programmatic access.
package export

import (
	"fmt"
	"io"
	"strings"
)

// RuleRow holds one rule's data for CSV export.
type RuleRow struct {
	Antecedent string
	Consequent string
	Support    float64
	Confidence float64
	Lift       float64
	Leverage   float64
	Conviction float64
}

// WriteRulesCSV writes rules as a CSV with header.
func WriteRulesCSV(w io.Writer, rules []RuleRow) error {
	header := "antecedent,consequent,support,confidence,lift,leverage,conviction\n"
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	for _, r := range rules {
		line := fmt.Sprintf("%s,%s,%.6f,%.6f,%.6f,%.6f,%.6f\n",
			quoteCSV(r.Antecedent), quoteCSV(r.Consequent),
			r.Support, r.Confidence, r.Lift, r.Leverage, r.Conviction)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

// FreqSetRow holds one frequent itemset for CSV export.
type FreqSetRow struct {
	Items   string
	Count   int
	Support float64
	Size    int
}

// WriteFreqSetsCSV writes frequent itemsets as CSV.
func WriteFreqSetsCSV(w io.Writer, sets []FreqSetRow) error {
	header := "items,count,support,size\n"
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	for _, s := range sets {
		line := fmt.Sprintf("%s,%d,%.6f,%d\n",
			quoteCSV(s.Items), s.Count, s.Support, s.Size)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

func quoteCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

// WriteSummary writes a human-readable text summary.
func WriteSummary(w io.Writer, transactions, freqSets, rules int, minSup, minConf float64) error {
	lines := []string{
		fmt.Sprintf("Mining Summary"),
		fmt.Sprintf("  Transactions:    %d", transactions),
		fmt.Sprintf("  Min Support:     %.4f", minSup),
		fmt.Sprintf("  Min Confidence:  %.4f", minConf),
		fmt.Sprintf("  Frequent Sets:   %d", freqSets),
		fmt.Sprintf("  Rules:           %d", rules),
		"",
	}
	for _, l := range lines {
		if _, err := io.WriteString(w, l+"\n"); err != nil {
			return err
		}
	}
	return nil
}
