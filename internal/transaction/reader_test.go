package transaction

import (
	"strings"
	"testing"
)

func TestReadText(t *testing.T) {
	input := "bread milk eggs\nbutter cheese\n# comment\n\napple banana\n"
	txs, err := ReadText(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadText: %v", err)
	}
	if len(txs) != 3 {
		t.Fatalf("got %d transactions, want 3", len(txs))
	}
	if txs[0][0] != "bread" {
		t.Errorf("first item = %q, want bread", txs[0][0])
	}
}

func TestReadJSON(t *testing.T) {
	input := `[["bread","milk"],["eggs","butter"]]`
	txs, err := ReadJSON(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("got %d transactions, want 2", len(txs))
	}
}

func TestVocabularyEncode(t *testing.T) {
	txs := []Transaction{{"bread", "milk"}, {"eggs", "bread"}}
	v := NewVocabulary(txs)
	if v.Size() != 3 {
		t.Fatalf("vocab size = %d, want 3", v.Size())
	}
	encoded := v.EncodeTransaction(txs[0])
	if len(encoded) != 2 {
		t.Fatalf("encoded len = %d, want 2", len(encoded))
	}
	decoded := v.DecodeTransaction(encoded)
	if len(decoded) != 2 {
		t.Fatalf("decoded len = %d, want 2", len(decoded))
	}
}
