// Package persist provides durable storage of mining results.
//
// Results are written as JSON with a trailing SHA-256 checksum line. On load,
// the checksum is verified to detect corruption or tampering. This ensures that
// saved frequent itemsets and rules can be reliably recovered and compared
// across runs without re-mining.
//
// File format:
//
//	{
//	  "version": 1,
//	  "total_transactions": N,
//	  "min_support": 0.2,
//	  "min_confidence": 0.6,
//	  "frequent_sets": [...],
//	  "rules": [...]
//	}
//	CHECKSUM:sha256:<hex>
package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrCorrupt indicates checksum mismatch on load.
	ErrCorrupt = errors.New("persist: checksum mismatch")

	// ErrNoChecksum indicates the file lacks a checksum trailer.
	ErrNoChecksum = errors.New("persist: missing checksum line")

	// ErrBadVersion indicates an unsupported file version.
	ErrBadVersion = errors.New("persist: unsupported version")
)

const currentVersion = 1

// FrequentSetRecord is the serializable form of a frequent itemset.
type FrequentSetRecord struct {
	Items   []string `json:"items"`
	Count   int      `json:"count"`
	Support float64  `json:"support"`
}

// RuleRecord is the serializable form of an association rule.
type RuleRecord struct {
	Antecedent []string `json:"antecedent"`
	Consequent []string `json:"consequent"`
	Support    float64  `json:"support"`
	Confidence float64  `json:"confidence"`
	Lift       float64  `json:"lift"`
	Leverage   float64  `json:"leverage"`
	Conviction float64  `json:"conviction"`
}

// MiningResult holds the complete output of a mining run.
type MiningResult struct {
	Version           int                 `json:"version"`
	TotalTransactions int                 `json:"total_transactions"`
	MinSupport        float64             `json:"min_support"`
	MinConfidence     float64             `json:"min_confidence"`
	FrequentSets      []FrequentSetRecord `json:"frequent_sets"`
	Rules             []RuleRecord        `json:"rules"`
}

// Save writes a mining result to path with an appended SHA-256 checksum.
func Save(path string, result *MiningResult) error {
	result.Version = currentVersion

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("persist: marshal: %w", err)
	}

	sum := sha256.Sum256(data)
	hexSum := hex.EncodeToString(sum[:])
	trailer := fmt.Sprintf("\nCHECKSUM:sha256:%s\n", hexSum)

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("persist: create: %w", err)
	}
	defer func() {
		f.Close()
		os.Remove(tmp)
	}()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("persist: write json: %w", err)
	}
	if _, err := f.WriteString(trailer); err != nil {
		return fmt.Errorf("persist: write checksum: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("persist: sync: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("persist: close: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("persist: rename: %w", err)
	}
	return nil
}

// Load reads a mining result and verifies its checksum.
func Load(path string) (*MiningResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("persist: read: %w", err)
	}

	content := string(raw)
	jsonData, checksum, err := splitChecksum(content)
	if err != nil {
		return nil, err
	}

	// verify checksum
	sum := sha256.Sum256([]byte(jsonData))
	computed := hex.EncodeToString(sum[:])
	if computed != checksum {
		return nil, ErrCorrupt
	}

	var result MiningResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, fmt.Errorf("persist: unmarshal: %w", err)
	}
	if result.Version != currentVersion {
		return nil, fmt.Errorf("%w: got %d", ErrBadVersion, result.Version)
	}
	return &result, nil
}

// splitChecksum separates the JSON payload from the checksum trailer.
func splitChecksum(content string) (jsonData, checksum string, err error) {
	// find last line starting with CHECKSUM:
	lines := strings.Split(content, "\n")
	checksumIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.HasPrefix(lines[i], "CHECKSUM:sha256:") {
			checksumIdx = i
			break
		}
	}
	if checksumIdx < 0 {
		return "", "", ErrNoChecksum
	}

	checksumLine := lines[checksumIdx]
	checksum = strings.TrimPrefix(checksumLine, "CHECKSUM:sha256:")
	checksum = strings.TrimSpace(checksum)

	// JSON is everything before the newline before checksum line
	jsonPart := strings.Join(lines[:checksumIdx], "\n")
	// trim trailing newline that was added before CHECKSUM
	jsonData = strings.TrimSuffix(jsonPart, "\n")

	return jsonData, checksum, nil
}

// Compare checks whether two mining results are logically equivalent
// (same frequent sets and rules, ignoring order).
func Compare(a, b *MiningResult) bool {
	if a.TotalTransactions != b.TotalTransactions {
		return false
	}
	if a.MinSupport != b.MinSupport || a.MinConfidence != b.MinConfidence {
		return false
	}
	if len(a.FrequentSets) != len(b.FrequentSets) {
		return false
	}
	if len(a.Rules) != len(b.Rules) {
		return false
	}

	// compare frequent sets by key
	aFS := make(map[string]int)
	for _, fs := range a.FrequentSets {
		aFS[itemsetKey(fs.Items)] = fs.Count
	}
	for _, fs := range b.FrequentSets {
		if aFS[itemsetKey(fs.Items)] != fs.Count {
			return false
		}
	}

	// compare rules by key
	aRules := make(map[string]bool)
	for _, r := range a.Rules {
		aRules[ruleKey(r)] = true
	}
	for _, r := range b.Rules {
		if !aRules[ruleKey(r)] {
			return false
		}
	}

	return true
}

func itemsetKey(items []string) string {
	return strings.Join(items, "\x00")
}

func ruleKey(r RuleRecord) string {
	return strings.Join(r.Antecedent, ",") + "=>" + strings.Join(r.Consequent, ",")
}
