// Package rank provides sorting and ranking of association rules by
// multiple interestingness criteria. Rules can be ranked by a single
// metric or by a weighted composite score.
package rank

import "sort"

// Rule represents a rule with its computed metrics for ranking purposes.
type Rule struct {
	Antecedent []string
	Consequent []string
	Support    float64
	Confidence float64
	Lift       float64
	Leverage   float64
	Conviction float64
	Cosine     float64
	Jaccard    float64
}

// Metric identifies which interestingness measure to rank by.
type Metric int

const (
	BySupport    Metric = iota
	ByConfidence
	ByLift
	ByLeverage
	ByConviction
	ByCosine
	ByJaccard
)

// SortBy sorts rules in descending order by the given metric.
func SortBy(rules []Rule, m Metric) {
	sort.SliceStable(rules, func(i, j int) bool {
		return metricValue(rules[i], m) > metricValue(rules[j], m)
	})
}

func metricValue(r Rule, m Metric) float64 {
	switch m {
	case BySupport:
		return r.Support
	case ByConfidence:
		return r.Confidence
	case ByLift:
		return r.Lift
	case ByLeverage:
		return r.Leverage
	case ByConviction:
		return r.Conviction
	case ByCosine:
		return r.Cosine
	case ByJaccard:
		return r.Jaccard
	default:
		return r.Confidence
	}
}

// TopK returns the top K rules by the given metric.
func TopK(rules []Rule, k int, m Metric) []Rule {
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	SortBy(sorted, m)
	if k > len(sorted) {
		k = len(sorted)
	}
	return sorted[:k]
}

// WeightedScore computes a composite score from multiple metrics.
type WeightedScore struct {
	SupportWeight    float64
	ConfidenceWeight float64
	LiftWeight       float64
	LeverageWeight   float64
	ConvictionWeight float64
}

// DefaultWeights returns equal weights for all metrics.
func DefaultWeights() WeightedScore {
	return WeightedScore{
		SupportWeight:    0.2,
		ConfidenceWeight: 0.3,
		LiftWeight:       0.3,
		LeverageWeight:   0.1,
		ConvictionWeight: 0.1,
	}
}

// Score computes the weighted composite score for a rule.
func (ws *WeightedScore) Score(r Rule) float64 {
	return ws.SupportWeight*r.Support +
		ws.ConfidenceWeight*r.Confidence +
		ws.LiftWeight*normalize(r.Lift, 10) +
		ws.LeverageWeight*normalize(r.Leverage, 1) +
		ws.ConvictionWeight*normalize(r.Conviction, 10)
}

// normalize maps a value to [0,1] using v/max clamping.
func normalize(v, max float64) float64 {
	if max <= 0 {
		return 0
	}
	if v > max {
		return 1
	}
	if v < 0 {
		return 0
	}
	return v / max
}

// SortByScore sorts rules by weighted composite score descending.
func SortByScore(rules []Rule, ws *WeightedScore) {
	sort.SliceStable(rules, func(i, j int) bool {
		return ws.Score(rules[i]) > ws.Score(rules[j])
	})
}

// Ranked wraps a rule with its computed rank and score.
type Ranked struct {
	Rule  Rule
	Rank  int
	Score float64
}

// RankAll assigns ranks to all rules by weighted score.
func RankAll(rules []Rule, ws *WeightedScore) []Ranked {
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	SortByScore(sorted, ws)
	ranked := make([]Ranked, len(sorted))
	for i, r := range sorted {
		ranked[i] = Ranked{Rule: r, Rank: i + 1, Score: ws.Score(r)}
	}
	return ranked
}
