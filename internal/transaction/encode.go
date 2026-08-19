package transaction

import (
	"fmt"
	"sort"
	"strings"
)

// Vocabulary maps items to integer IDs for compact internal representation.
type Vocabulary struct {
	itemToID map[Item]int
	idToItem []Item
}

// NewVocabulary builds a vocabulary from all unique items in the transactions.
// Items are assigned IDs in alphabetical order for determinism.
func NewVocabulary(txs []Transaction) *Vocabulary {
	unique := make(map[Item]struct{})
	for _, tx := range txs {
		for _, it := range tx {
			unique[it] = struct{}{}
		}
	}
	items := make([]Item, 0, len(unique))
	for it := range unique {
		items = append(items, it)
	}
	sort.Strings(items)

	v := &Vocabulary{
		itemToID: make(map[Item]int, len(items)),
		idToItem: items,
	}
	for i, it := range items {
		v.itemToID[it] = i
	}
	return v
}

// Size returns the number of items in the vocabulary.
func (v *Vocabulary) Size() int { return len(v.idToItem) }

// Encode converts an item to its integer ID. Returns -1 if unknown.
func (v *Vocabulary) Encode(item Item) int {
	id, ok := v.itemToID[item]
	if !ok {
		return -1
	}
	return id
}

// Decode converts an integer ID back to the item string.
func (v *Vocabulary) Decode(id int) (Item, error) {
	if id < 0 || id >= len(v.idToItem) {
		return "", fmt.Errorf("id %d out of range [0, %d)", id, len(v.idToItem))
	}
	return v.idToItem[id], nil
}

// EncodeTransaction converts a transaction to a sorted slice of IDs.
func (v *Vocabulary) EncodeTransaction(tx Transaction) []int {
	ids := make([]int, 0, len(tx))
	for _, it := range tx {
		if id, ok := v.itemToID[it]; ok {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)
	return ids
}

// DecodeTransaction converts a slice of IDs back to items.
func (v *Vocabulary) DecodeTransaction(ids []int) Transaction {
	tx := make(Transaction, 0, len(ids))
	for _, id := range ids {
		if id >= 0 && id < len(v.idToItem) {
			tx = append(tx, v.idToItem[id])
		}
	}
	return tx
}

// EncodeAll encodes all transactions to integer representation.
func (v *Vocabulary) EncodeAll(txs []Transaction) [][]int {
	encoded := make([][]int, len(txs))
	for i, tx := range txs {
		encoded[i] = v.EncodeTransaction(tx)
	}
	return encoded
}

// String returns a human-readable vocabulary listing.
func (v *Vocabulary) String() string {
	var b strings.Builder
	for id, item := range v.idToItem {
		fmt.Fprintf(&b, "%d: %s\n", id, item)
	}
	return b.String()
}
