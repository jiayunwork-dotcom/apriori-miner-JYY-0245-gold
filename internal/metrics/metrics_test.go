package metrics

import (
	"math"
	"testing"
)

const eps = 1e-9

func approx(a, b float64) bool {
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	return math.Abs(a-b) < eps
}

func TestComputeLift(t *testing.T) {
	// P(XY)=0.4, P(X)=0.6, P(Y)=0.5 => lift = 0.4/(0.6*0.5) = 1.333
	p := Params{SupportXY: 0.4, SupportX: 0.6, SupportY: 0.5, Confidence: 0.4 / 0.6}
	m := Compute(p)
	want := 0.4 / (0.6 * 0.5)
	if !approx(m.Lift, want) {
		t.Errorf("Lift = %f, want %f", m.Lift, want)
	}
}

func TestComputeLiftZeroDenom(t *testing.T) {
	p := Params{SupportXY: 0.3, SupportX: 0.0, SupportY: 0.5}
	m := Compute(p)
	if m.Lift != 0 {
		t.Errorf("Lift with zero SupportX = %f, want 0", m.Lift)
	}
}

func TestComputeLeverage(t *testing.T) {
	p := Params{SupportXY: 0.4, SupportX: 0.6, SupportY: 0.5}
	m := Compute(p)
	want := 0.4 - 0.6*0.5
	if !approx(m.Leverage, want) {
		t.Errorf("Leverage = %f, want %f", m.Leverage, want)
	}
}

func TestComputeConviction(t *testing.T) {
	// confidence=0.8, P(Y)=0.5 => conviction = (1-0.5)/(1-0.8) = 2.5
	p := Params{SupportXY: 0.4, SupportX: 0.5, SupportY: 0.5, Confidence: 0.8}
	m := Compute(p)
	want := 0.5 / 0.2
	if !approx(m.Conviction, want) {
		t.Errorf("Conviction = %f, want %f", m.Conviction, want)
	}
}

func TestComputeConvictionPerfectConfidence(t *testing.T) {
	p := Params{SupportXY: 0.5, SupportX: 0.5, SupportY: 0.5, Confidence: 1.0}
	m := Compute(p)
	if !math.IsInf(m.Conviction, 1) {
		t.Errorf("Conviction with confidence=1 = %f, want +Inf", m.Conviction)
	}
}

func TestComputeCosine(t *testing.T) {
	// cosine = P(XY)/sqrt(P(X)*P(Y)) = 0.4/sqrt(0.6*0.5)
	p := Params{SupportXY: 0.4, SupportX: 0.6, SupportY: 0.5}
	m := Compute(p)
	want := 0.4 / math.Sqrt(0.6*0.5)
	if !approx(m.Cosine, want) {
		t.Errorf("Cosine = %f, want %f", m.Cosine, want)
	}
}

func TestComputeJaccard(t *testing.T) {
	// countXY=20, countX=30, countY=40 => union=30+40-20=50, J=20/50=0.4
	p := Params{CountXY: 20, CountX: 30, CountY: 40, Total: 100}
	m := Compute(p)
	want := 20.0 / 50.0
	if !approx(m.Jaccard, want) {
		t.Errorf("Jaccard = %f, want %f", m.Jaccard, want)
	}
}

func TestComputeJaccardZero(t *testing.T) {
	p := Params{CountXY: 0, CountX: 0, CountY: 0}
	m := Compute(p)
	if m.Jaccard != 0 {
		t.Errorf("Jaccard zero case = %f, want 0", m.Jaccard)
	}
}

func TestComputeKulczynski(t *testing.T) {
	// P(Y|X) = P(XY)/P(X) = 0.4/0.6; P(X|Y) = P(XY)/P(Y) = 0.4/0.5
	// Kulc = avg = (0.667 + 0.8)/2 = 0.733
	p := Params{SupportXY: 0.4, SupportX: 0.6, SupportY: 0.5}
	m := Compute(p)
	want := ((0.4 / 0.6) + (0.4 / 0.5)) / 2
	if !approx(m.Kulczynski, want) {
		t.Errorf("Kulczynski = %f, want %f", m.Kulczynski, want)
	}
}

func TestBatchCompute(t *testing.T) {
	params := []Params{
		{SupportXY: 0.3, SupportX: 0.5, SupportY: 0.6, Confidence: 0.6, CountXY: 30, CountX: 50, CountY: 60, Total: 100},
		{SupportXY: 0.2, SupportX: 0.4, SupportY: 0.5, Confidence: 0.5, CountXY: 20, CountX: 40, CountY: 50, Total: 100},
	}
	results := BatchCompute(params)
	if len(results) != 2 {
		t.Fatalf("BatchCompute returned %d results, want 2", len(results))
	}
	// verify first result lift
	want := 0.3 / (0.5 * 0.6)
	if !approx(results[0].Lift, want) {
		t.Errorf("results[0].Lift = %f, want %f", results[0].Lift, want)
	}
}

func TestIsPositivelyCorrelated(t *testing.T) {
	m := RuleMetrics{Lift: 1.5}
	if !IsPositivelyCorrelated(m) {
		t.Error("Lift=1.5 should be positively correlated")
	}
	m.Lift = 0.8
	if IsPositivelyCorrelated(m) {
		t.Error("Lift=0.8 should not be positively correlated")
	}
}

func TestIsNegativelyCorrelated(t *testing.T) {
	m := RuleMetrics{Lift: 0.7}
	if !IsNegativelyCorrelated(m) {
		t.Error("Lift=0.7 should be negatively correlated")
	}
	m.Lift = 1.2
	if IsNegativelyCorrelated(m) {
		t.Error("Lift=1.2 should not be negatively correlated")
	}
}

func TestIsIndependent(t *testing.T) {
	m := RuleMetrics{Lift: 1.0001}
	if !IsIndependent(m, 0.001) {
		t.Error("Lift=1.0001 should be independent within epsilon=0.001")
	}
	if IsIndependent(m, 0.00001) {
		t.Error("Lift=1.0001 should not be independent within epsilon=0.00001")
	}
}

func TestLeverageSign(t *testing.T) {
	// positive leverage when items co-occur more than expected
	pos := Params{SupportXY: 0.5, SupportX: 0.6, SupportY: 0.7}
	m := Compute(pos)
	if m.Leverage <= 0 {
		t.Errorf("expected positive leverage, got %f", m.Leverage)
	}

	// negative leverage
	neg := Params{SupportXY: 0.1, SupportX: 0.6, SupportY: 0.7}
	m = Compute(neg)
	if m.Leverage >= 0 {
		t.Errorf("expected negative leverage, got %f", m.Leverage)
	}
}
