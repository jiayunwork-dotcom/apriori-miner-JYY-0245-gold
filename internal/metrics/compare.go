package metrics

import "math"

// Comparison holds the result of comparing two rules' metrics.
type Comparison struct {
	LiftDiff       float64
	LeverageDiff   float64
	ConvictionDiff float64
	CosineDiff     float64
	JaccardDiff    float64
	Dominant       string // "first", "second", or "mixed"
}

// Compare computes the difference between two sets of rule metrics.
func Compare(a, b RuleMetrics) Comparison {
	c := Comparison{
		LiftDiff:       a.Lift - b.Lift,
		LeverageDiff:   a.Leverage - b.Leverage,
		ConvictionDiff: convDiff(a.Conviction, b.Conviction),
		CosineDiff:     a.Cosine - b.Cosine,
		JaccardDiff:    a.Jaccard - b.Jaccard,
	}
	pos, neg := 0, 0
	diffs := []float64{c.LiftDiff, c.LeverageDiff, c.CosineDiff, c.JaccardDiff}
	for _, d := range diffs {
		if d > 0 {
			pos++
		} else if d < 0 {
			neg++
		}
	}
	if pos > neg {
		c.Dominant = "first"
	} else if neg > pos {
		c.Dominant = "second"
	} else {
		c.Dominant = "mixed"
	}
	return c
}

func convDiff(a, b float64) float64 {
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return 0
	}
	if math.IsInf(a, 1) {
		return 1000
	}
	if math.IsInf(b, 1) {
		return -1000
	}
	return a - b
}

// Dominates returns true if every metric of a is >= every metric of b
// (Pareto dominance).
func Dominates(a, b RuleMetrics) bool {
	if a.Lift < b.Lift {
		return false
	}
	if a.Leverage < b.Leverage {
		return false
	}
	if a.Cosine < b.Cosine {
		return false
	}
	if a.Jaccard < b.Jaccard {
		return false
	}
	if a.Kulczynski < b.Kulczynski {
		return false
	}
	// At least one must be strictly greater
	return a.Lift > b.Lift || a.Leverage > b.Leverage ||
		a.Cosine > b.Cosine || a.Jaccard > b.Jaccard ||
		a.Kulczynski > b.Kulczynski
}

// ParetoFront returns the non-dominated rules from a set of metrics.
// A rule is in the Pareto front if no other rule dominates it.
func ParetoFront(metrics []RuleMetrics) []int {
	var front []int
	for i := range metrics {
		dominated := false
		for j := range metrics {
			if i != j && Dominates(metrics[j], metrics[i]) {
				dominated = true
				break
			}
		}
		if !dominated {
			front = append(front, i)
		}
	}
	return front
}

// EuclideanDistance computes the Euclidean distance between two rule
// metrics in the normalized metric space.
func EuclideanDistance(a, b RuleMetrics) float64 {
	d := 0.0
	d += (a.Lift - b.Lift) * (a.Lift - b.Lift)
	d += (a.Leverage - b.Leverage) * (a.Leverage - b.Leverage)
	d += (a.Cosine - b.Cosine) * (a.Cosine - b.Cosine)
	d += (a.Jaccard - b.Jaccard) * (a.Jaccard - b.Jaccard)
	d += (a.Kulczynski - b.Kulczynski) * (a.Kulczynski - b.Kulczynski)
	return math.Sqrt(d)
}
