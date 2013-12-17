package main

import (
	"math"
	"testing"
)

func TestAnalyze(t *testing.T) {
	points := map[int]float64{
		1:   65,
		18:  996,
		36:  1652,
		72:  1853,
		108: 1829,
		144: 1775,
		216: 1702,
	}

	m, err := analyze(points)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(m.alpha-0.0203030740304324) > 0.00001 {
		t.Errorf("Bad alpha: %f", m.alpha)
	}

	if math.Abs(m.beta-6.7246130982513e-5) > 0.00001 {
		t.Errorf("Bad beta: %f", m.beta)
	}

	if m.y != 65 {
		t.Errorf("Bad Y: %f", m.y)
	}

	if m.nmax != 120 {
		t.Errorf("Bad Nmax: %d", m.nmax)
	}

	if m.nopt != 50 {
		t.Errorf("Bad Nopt: %d", m.nopt)
	}
}
