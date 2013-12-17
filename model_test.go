package usl

import (
	"math"
	"testing"
)

func TestAnalyze(t *testing.T) {
	measurements := MeasurementSet{
		Measurement{X: 1, Y: 65},
		Measurement{X: 18, Y: 996},
		Measurement{X: 36, Y: 1652},
		Measurement{X: 72, Y: 1853},
		Measurement{X: 108, Y: 1829},
		Measurement{X: 144, Y: 1775},
		Measurement{X: 216, Y: 1702},
	}

	m, err := Build(measurements)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(m.Alpha-0.0203030740304324) > 0.00001 {
		t.Errorf("Bad alpha: %f", m.Alpha)
	}

	if math.Abs(m.Beta-6.7246130982513e-5) > 0.00001 {
		t.Errorf("Bad beta: %f", m.Beta)
	}

	if m.Y != 65 {
		t.Errorf("Bad Y: %f", m.Y)
	}

	if m.Nmax != 120 {
		t.Errorf("Bad Nmax: %d", m.Nmax)
	}

	expected := 1164.4929146148988
	actual := m.Predict(500.0)
	if math.Abs(expected-actual) > 0.00001 {
		t.Errorf("Expected %v but was %v", expected, actual)
	}
}
