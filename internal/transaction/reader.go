// Package transaction provides parsing and encoding of transaction data in
// multiple formats: whitespace-separated text (one basket per line), CSV,
// and JSON arrays. It decouples I/O from the mining algorithm so the core
// Apriori package remains format-agnostic.
package transaction

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Item is a single product identifier within a basket.
type Item = string

// Transaction is one basket of items.
type Transaction = []Item

// ReadText reads transactions from whitespace-separated lines. Empty lines
// and lines starting with '#' are skipped.
func ReadText(r io.Reader) ([]Transaction, error) {
	var txs []Transaction
	sc := bufio.NewScanner(r)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items := strings.Fields(line)
		if len(items) == 0 {
			continue
		}
		txs = append(txs, items)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("transaction: read line %d: %w", lineNo, err)
	}
	return txs, nil
}

// ReadCSV reads transactions from CSV where each row is one basket and
// each field is one item. Empty fields are skipped.
func ReadCSV(r io.Reader) ([]Transaction, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1 // variable
	cr.TrimLeadingSpace = true

	var txs []Transaction
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("transaction: csv: %w", err)
		}
		var items []Item
		for _, f := range row {
			f = strings.TrimSpace(f)
			if f != "" {
				items = append(items, f)
			}
		}
		if len(items) > 0 {
			txs = append(txs, items)
		}
	}
	return txs, nil
}

// ReadJSON reads transactions from a JSON array of arrays:
// [["bread","milk"],["eggs","butter"]].
func ReadJSON(r io.Reader) ([]Transaction, error) {
	var raw [][]string
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("transaction: json: %w", err)
	}
	txs := make([]Transaction, 0, len(raw))
	for _, items := range raw {
		if len(items) > 0 {
			txs = append(txs, items)
		}
	}
	return txs, nil
}

// WriteText writes transactions as whitespace-separated lines.
func WriteText(w io.Writer, txs []Transaction) error {
	bw := bufio.NewWriter(w)
	for _, tx := range txs {
		if _, err := bw.WriteString(strings.Join(tx, " ")); err != nil {
			return err
		}
		if err := bw.WriteByte('\n'); err != nil {
			return err
		}
	}
	return bw.Flush()
}

// WriteJSON writes transactions as a JSON array of arrays.
func WriteJSON(w io.Writer, txs []Transaction) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(txs)
}
