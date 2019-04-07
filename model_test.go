package usl

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

func TestMeasurementSorting(t *testing.T) {
	measurements := MeasurementSet{
		Measurement{X: 72, Y: 1853},
		Measurement{X: 108, Y: 1829},
		Measurement{X: 18, Y: 996},
		Measurement{X: 216, Y: 1702},
		Measurement{X: 36, Y: 1652},
		Measurement{X: 144, Y: 1775},
		Measurement{X: 1, Y: 65},
	}

	sort.Sort(measurements)

	expected := "1,18,36,72,108,144,216,"
	actual := ""
	for _, m := range measurements {
		actual += fmt.Sprintf("%v,", m.X)
	}

	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func TestBuild(t *testing.T) {
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
		t.Errorf("Bad alpha: %v", m.Alpha)
	}

	if math.Abs(m.Beta-6.7246130982513e-5) > 0.00001 {
		t.Errorf("Bad beta: %v", m.Beta)
	}

	if m.Y != 65 {
		t.Errorf("Bad Y: %v", m.Y)
	}

	if m.Peak != 120 {
		t.Errorf("Bad Peak: %v", m.Peak)
	}

	expected := 1164.4929146148988
	actual := m.Predict(500.0)
	if math.Abs(expected-actual) > 0.00001 {
		t.Errorf("Expected %v but was %v", expected, actual)
	}
}
