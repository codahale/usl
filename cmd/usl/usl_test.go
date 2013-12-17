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

func TestBadLine(t *testing.T) {
	m, err := parseLine(0, []string{"funk"})
	if err == nil {
		t.Fatalf("Shouldn't have parsed, but returned %v", m)
	}

	expected := "invalid line at line 1"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func TestBadX(t *testing.T) {
	m, err := parseLine(0, []string{"f", "1"})
	if err == nil {
		t.Fatalf("Shouldn't have parsed, but returned %v", m)
	}

	expected := "strconv.ParseFloat: parsing \"f\": invalid syntax at line 1, column 1"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func TestBadY(t *testing.T) {
	m, err := parseLine(0, []string{"1", "f"})
	if err == nil {
		t.Fatalf("Shouldn't have parsed, but returned %v", m)
	}

	expected := "strconv.ParseFloat: parsing \"f\": invalid syntax at line 1, column 2"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}
