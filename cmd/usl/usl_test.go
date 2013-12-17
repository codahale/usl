package main

import (
	"testing"

	"github.com/codahale/usl"
)

func TestParsing(t *testing.T) {
	expected := usl.MeasurementSet{
		usl.Measurement{X: 1, Y: 65},
		usl.Measurement{X: 18, Y: 996},
		usl.Measurement{X: 36, Y: 1652},
		usl.Measurement{X: 72, Y: 1853},
		usl.Measurement{X: 108, Y: 1829},
		usl.Measurement{X: 144, Y: 1775},
		usl.Measurement{X: 216, Y: 1702},
	}

	actual, err := parseCSV("example.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("Expected %d measurements, but was %d", len(expected), len(actual))
	}

	for i, a := range actual {
		e := expected[i]

		if a.X != e.X || a.Y != e.Y {
			t.Fatalf("Expected %v, but was %v", e, a)
		}
	}
}
