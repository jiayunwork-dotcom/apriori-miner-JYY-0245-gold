package sampling

import "testing"

func TestRandomSample(t *testing.T) {
	txs := []Transaction{
		{"a", "b"}, {"c", "d"}, {"e", "f"}, {"g", "h"},
	}
	sample := RandomSample(txs, 2, 42)
	if len(sample) != 2 {
		t.Fatalf("sample size = %d, want 2", len(sample))
	}
}

func TestRandomSampleWithoutReplacement(t *testing.T) {
	txs := []Transaction{
		{"a"}, {"b"}, {"c"}, {"d"}, {"e"},
	}
	sample := RandomSampleWithoutReplacement(txs, 3, 99)
	if len(sample) != 3 {
		t.Fatalf("sample size = %d, want 3", len(sample))
	}
	// Check no duplicates (by checking they are from original)
	seen := make(map[string]bool)
	for _, tx := range sample {
		k := tx[0]
		if seen[k] {
			t.Errorf("duplicate transaction: %v", tx)
		}
		seen[k] = true
	}
}

func TestBootstrapSample(t *testing.T) {
	txs := []Transaction{{"x", "y"}, {"z", "w"}}
	bs := BootstrapSample(txs, 7)
	if len(bs) != 2 {
		t.Fatalf("bootstrap size = %d, want 2 (same as input)", len(bs))
	}
}
