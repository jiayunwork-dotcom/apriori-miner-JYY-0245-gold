// Package metrics computes interestingness measures for association rules.
//
// Given support values for antecedent, consequent, and the union, it computes:
//   - Lift: P(X∪Y) / (P(X)*P(Y)). Lift > 1 means positive correlation.
//   - Leverage: P(X∪Y) - P(X)*P(Y). How much more co-occurrence than expected.
//   - Conviction: (1-P(Y)) / (1-Confidence). Infinity when confidence=1.
//   - Cosine: P(X∪Y) / sqrt(P(X)*P(Y)). Geometric mean normalized co-occurrence.
//   - Jaccard: |X∩Y| / |X∪Y|, approximated via support counts.
//   - Kulczynski: average of P(Y|X) and P(X|Y).
package metrics

import "math"

// RuleMetrics holds all interestingness measures for a single rule.
type RuleMetrics struct {
	Lift       float64
	Leverage   float64
	Conviction float64
	Cosine     float64
	Jaccard    float64
	Kulczynski float64
}

// Params are the inputs needed to compute rule metrics.
type Params struct {
	// SupportXY is P(X ∪ Y) = count(X∪Y) / total.
	SupportXY float64
	// SupportX is P(X) = count(X) / total.
	SupportX float64
	// SupportY is P(Y) = count(Y) / total.
	SupportY float64
	// Confidence is P(Y|X) = P(X∪Y)/P(X).
	Confidence float64
	// CountXY is the absolute support of the union.
	CountXY int
	// CountX is the absolute support of the antecedent.
	CountX int
	// CountY is the absolute support of the consequent.
	CountY int
	// Total is the number of transactions.
	Total int
}

// Compute calculates all interestingness measures from the given parameters.
func Compute(p Params) RuleMetrics {
	var m RuleMetrics

	m.Lift = computeLift(p.SupportXY, p.SupportX, p.SupportY)
	m.Leverage = computeLeverage(p.SupportXY, p.SupportX, p.SupportY)
	m.Conviction = computeConviction(p.SupportY, p.Confidence)
	m.Cosine = computeCosine(p.SupportXY, p.SupportX, p.SupportY)
	m.Jaccard = computeJaccard(p.CountXY, p.CountX, p.CountY)
	m.Kulczynski = computeKulczynski(p.SupportXY, p.SupportX, p.SupportY)

	return m
}

func computeLift(supXY, supX, supY float64) float64 {
	denom := supX * supY
	if denom == 0 {
		return 0
	}
	return supXY / denom
}

func computeLeverage(supXY, supX, supY float64) float64 {
	return supXY - supX*supY
}

func computeConviction(supY, confidence float64) float64 {
	denom := 1 - confidence
	if denom <= 0 {
		return math.Inf(1)
	}
	num := 1 - supY
	if num <= 0 {
		return 0
	}
	return num / denom
}

func computeCosine(supXY, supX, supY float64) float64 {
	denom := math.Sqrt(supX * supY)
	if denom == 0 {
		return 0
	}
	return supXY / denom
}

func computeJaccard(countXY, countX, countY int) float64 {
	union := countX + countY - countXY
	if union == 0 {
		return 0
	}
	return float64(countXY) / float64(union)
}

func computeKulczynski(supXY, supX, supY float64) float64 {
	var sum float64
	n := 0
	if supX > 0 {
		sum += supXY / supX
		n++
	}
	if supY > 0 {
		sum += supXY / supY
		n++
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

// BatchCompute calculates metrics for a batch of rules specified as params.
func BatchCompute(params []Params) []RuleMetrics {
	results := make([]RuleMetrics, len(params))
	for i, p := range params {
		results[i] = Compute(p)
	}
	return results
}

// IsPositivelyCorrelated returns true if lift > 1 (items co-occur more than
// expected by chance).
func IsPositivelyCorrelated(m RuleMetrics) bool {
	return m.Lift > 1.0
}

// IsNegativelyCorrelated returns true if lift < 1.
func IsNegativelyCorrelated(m RuleMetrics) bool {
	return m.Lift < 1.0
}

// IsIndependent returns true if lift is approximately 1 (within epsilon).
func IsIndependent(m RuleMetrics, epsilon float64) bool {
	return math.Abs(m.Lift-1.0) <= epsilon
}
