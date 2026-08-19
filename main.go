// Command apriori-miner mines frequent itemsets and association rules from
// transaction data (one basket per line, items separated by whitespace).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"apriori-miner/internal/apriori"
)

func main() {
	input := flag.String("input", "transactions.txt", "path to transactions file")
	minSup := flag.Float64("min-support", 0.2, "minimum support as a fraction [0,1]")
	minConf := flag.Float64("min-confidence", 0.6, "minimum rule confidence [0,1]")
	top := flag.Int("top", 20, "max number of rules to print")
	flag.Parse()

	txs, err := readTransactions(*input)
	if err != nil {
		fatal("read %s: %v", *input, err)
	}
	if len(txs) == 0 {
		fatal("no transactions read from %s", *input)
	}

	freq := apriori.Apriori(txs, *minSup)
	rules := apriori.GenerateRules(freq, len(txs), *minConf)

	fmt.Printf("transactions: %d, frequent itemsets: %d, rules: %d\n\n",
		len(txs), len(freq), len(rules))

	sortedFreq := append([]apriori.FrequentSet(nil), freq...)
	sort.Slice(sortedFreq, func(i, j int) bool {
		if sortedFreq[i].Count != sortedFreq[j].Count {
			return sortedFreq[i].Count > sortedFreq[j].Count
		}
		return keyStr(sortedFreq[i].Items) < keyStr(sortedFreq[j].Items)
	})
	fmt.Println("Frequent itemsets:")
	for _, f := range sortedFreq {
		fmt.Printf("  [%d] %s (support %.2f)\n", f.Count, joinItems(f.Items), float64(f.Count)/float64(len(txs)))
	}

	fmt.Println("\nTop rules (confidence desc):")
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Confidence > rules[j].Confidence
	})
	limit := *top
	if limit > len(rules) {
		limit = len(rules)
	}
	for _, r := range rules[:limit] {
		fmt.Printf("  %s => %s  (sup %.2f, conf %.2f)\n",
			joinItems(r.Antecedent), joinItems(r.Consequent), r.Support, r.Confidence)
	}
}

func readTransactions(path string) ([]apriori.Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var txs []apriori.Transaction
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var t apriori.Transaction
		for _, it := range strings.Fields(line) {
			t = append(t, apriori.Item(it))
		}
		txs = append(txs, t)
	}
	return txs, sc.Err()
}

func joinItems(s apriori.Itemset) string {
	parts := make([]string, len(s))
	for i, it := range s {
		parts[i] = string(it)
	}
	return strings.Join(parts, ",")
}

func keyStr(s apriori.Itemset) string {
	parts := make([]string, len(s))
	for i, it := range s {
		parts[i] = string(it)
	}
	return strings.Join(parts, "\x00")
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "apriori-miner: "+format+"\n", args...)
	os.Exit(1)
}
