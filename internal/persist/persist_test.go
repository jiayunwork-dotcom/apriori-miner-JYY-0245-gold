package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func sampleResult() *MiningResult {
	return &MiningResult{
		TotalTransactions: 100,
		MinSupport:        0.2,
		MinConfidence:     0.6,
		FrequentSets: []FrequentSetRecord{
			{Items: []string{"bread"}, Count: 60, Support: 0.6},
			{Items: []string{"milk"}, Count: 50, Support: 0.5},
			{Items: []string{"bread", "milk"}, Count: 30, Support: 0.3},
		},
		Rules: []RuleRecord{
			{
				Antecedent: []string{"bread"},
				Consequent: []string{"milk"},
				Support:    0.3,
				Confidence: 0.5,
				Lift:       1.0,
				Leverage:   0.0,
				Conviction: 1.0,
			},
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "result.json")

	original := sampleResult()
	if err := Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.TotalTransactions != original.TotalTransactions {
		t.Errorf("TotalTransactions = %d, want %d", loaded.TotalTransactions, original.TotalTransactions)
	}
	if loaded.MinSupport != original.MinSupport {
		t.Errorf("MinSupport = %f, want %f", loaded.MinSupport, original.MinSupport)
	}
	if len(loaded.FrequentSets) != len(original.FrequentSets) {
		t.Errorf("FrequentSets count = %d, want %d", len(loaded.FrequentSets), len(original.FrequentSets))
	}
	if len(loaded.Rules) != len(original.Rules) {
		t.Errorf("Rules count = %d, want %d", len(loaded.Rules), len(original.Rules))
	}
}

func TestLoadCorruptChecksum(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")

	if err := Save(path, sampleResult()); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	// flip a byte in the JSON payload
	modified := make([]byte, len(data))
	copy(modified, data)
	modified[10] ^= 0xFF
	os.WriteFile(path, modified, 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestLoadMissingChecksum(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nochecksum.json")

	os.WriteFile(path, []byte(`{"version":1,"total_transactions":5}`), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected missing checksum error")
	}
}

func TestLoadBadVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "badver.json")

	result := sampleResult()
	if err := Save(path, result); err != nil {
		t.Fatal(err)
	}

	// manually re-write with version=99, recompute checksum
	result.Version = 99
	// write without using Save (which forces version=1)
	writeWithChecksum(t, path, result)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected version error")
	}
}

func TestCompareEqual(t *testing.T) {
	a := sampleResult()
	b := sampleResult()
	if !Compare(a, b) {
		t.Error("identical results should compare equal")
	}
}

func TestCompareDifferentCounts(t *testing.T) {
	a := sampleResult()
	b := sampleResult()
	b.FrequentSets[0].Count = 999
	if Compare(a, b) {
		t.Error("different counts should not compare equal")
	}
}

func TestCompareDifferentRules(t *testing.T) {
	a := sampleResult()
	b := sampleResult()
	b.Rules = append(b.Rules, RuleRecord{
		Antecedent: []string{"extra"},
		Consequent: []string{"rule"},
	})
	if Compare(a, b) {
		t.Error("different rule sets should not compare equal")
	}
}

func TestSaveAtomicity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "atomic.json")

	if err := Save(path, sampleResult()); err != nil {
		t.Fatal(err)
	}
	tmp := path + ".tmp"
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Error(".tmp file should not remain after successful save")
	}
}

func TestSaveLoadEmptyResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	result := &MiningResult{
		TotalTransactions: 0,
		MinSupport:        0.1,
		MinConfidence:     0.5,
		FrequentSets:      nil,
		Rules:             nil,
	}
	if err := Save(path, result); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.TotalTransactions != 0 {
		t.Errorf("TotalTransactions = %d, want 0", loaded.TotalTransactions)
	}
	if len(loaded.FrequentSets) != 0 {
		t.Errorf("FrequentSets = %v, want empty", loaded.FrequentSets)
	}
}

func TestSaveLoadLargeResult(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")

	result := &MiningResult{
		TotalTransactions: 10000,
		MinSupport:        0.01,
		MinConfidence:     0.3,
	}
	for i := 0; i < 200; i++ {
		result.FrequentSets = append(result.FrequentSets, FrequentSetRecord{
			Items:   []string{fmt.Sprintf("item_%d", i)},
			Count:   100 + i,
			Support: float64(100+i) / 10000.0,
		})
	}
	if err := Save(path, result); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.FrequentSets) != 200 {
		t.Errorf("FrequentSets = %d, want 200", len(loaded.FrequentSets))
	}
}

func TestRoundTripPreservesMetrics(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "metrics.json")

	result := &MiningResult{
		TotalTransactions: 50,
		MinSupport:        0.2,
		MinConfidence:     0.5,
		FrequentSets: []FrequentSetRecord{
			{Items: []string{"a", "b"}, Count: 15, Support: 0.3},
		},
		Rules: []RuleRecord{
			{
				Antecedent: []string{"a"},
				Consequent: []string{"b"},
				Support:    0.3,
				Confidence: 0.75,
				Lift:       1.5,
				Leverage:   0.1,
				Conviction: 2.0,
			},
		},
	}
	Save(path, result)
	loaded, _ := Load(path)

	r := loaded.Rules[0]
	if r.Lift != 1.5 {
		t.Errorf("Lift = %f, want 1.5", r.Lift)
	}
	if r.Leverage != 0.1 {
		t.Errorf("Leverage = %f, want 0.1", r.Leverage)
	}
	if r.Conviction != 2.0 {
		t.Errorf("Conviction = %f, want 2.0", r.Conviction)
	}
}

// helper to write with a custom version (bypassing Save's version enforcement)
func writeWithChecksum(t *testing.T, path string, result *MiningResult) {
	t.Helper()
	data, _ := json.MarshalIndent(result, "", "  ")
	sum := sha256.Sum256(data)
	hexSum := hex.EncodeToString(sum[:])
	content := string(data) + "\nCHECKSUM:sha256:" + hexSum + "\n"
	os.WriteFile(path, []byte(content), 0644)
}
