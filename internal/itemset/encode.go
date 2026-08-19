package itemset

import (
	"fmt"
	"strings"
)

// Key returns a stable string key for a set, suitable as a map key.
// Items are joined with a null byte separator.
func Key(s Set) string {
	return strings.Join(s, "\x00")
}

// FromKey reconstructs a Set from a key produced by Key().
func FromKey(k string) Set {
	if k == "" {
		return nil
	}
	return strings.Split(k, "\x00")
}

// Display returns a human-readable representation: {a, b, c}.
func Display(s Set) string {
	if len(s) == 0 {
		return "{}"
	}
	return "{" + strings.Join(s, ", ") + "}"
}

// DisplayRule formats an association rule as "A, B => C".
func DisplayRule(antecedent, consequent Set) string {
	return Display(antecedent) + " => " + Display(consequent)
}

// ParseDisplay parses a display-format set back into a Set.
// Input should be like "{a, b, c}" or "a, b, c".
func ParseDisplay(s string) Set {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var items []Item
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return New(items...)
}

// Encode converts a set into a compact binary representation where each
// item is length-prefixed. Format: [len1][item1][len2][item2]...
func Encode(s Set) []byte {
	var buf []byte
	for _, item := range s {
		b := []byte(item)
		if len(b) > 255 {
			b = b[:255]
		}
		buf = append(buf, byte(len(b)))
		buf = append(buf, b...)
	}
	return buf
}

// Decode reconstructs a Set from the binary format produced by Encode.
func Decode(data []byte) (Set, error) {
	var items []Item
	i := 0
	for i < len(data) {
		if i >= len(data) {
			return nil, fmt.Errorf("truncated at offset %d", i)
		}
		length := int(data[i])
		i++
		if i+length > len(data) {
			return nil, fmt.Errorf("item at offset %d extends past end", i)
		}
		items = append(items, string(data[i:i+length]))
		i += length
	}
	return items, nil
}
