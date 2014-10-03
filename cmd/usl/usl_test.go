package main

import (
	"os"
	"testing"

	"io/ioutil"
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

	actual, err := parseCSV("example.csv", 1, 2, false)
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
	m, err := parseLine(0, 1, 2, []string{"funk"})
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
	m, err := parseLine(0, 1, 2, []string{"f", "1"})
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
	m, err := parseLine(0, 1, 2, []string{"1", "f"})
	if err == nil {
		t.Fatalf("Shouldn't have parsed, but returned %v", m)
	}

	expected := "strconv.ParseFloat: parsing \"f\": invalid syntax at line 1, column 2"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func TestMainRun(t *testing.T) {
	expected := `1.000000,65.000000
2.000000,127.396329
3.000000,187.318153
`
	stdout, stderr := fakeMain(t, "-in", "example.csv", "1", "2", "3")

	actual := string(stdout)
	if expected != actual {
		t.Errorf("Expected\n%s\nbut was\n%s", expected, actual)
	}

	expected = `Model:
	α:    0.020303 (constrained by contention effects)
	β:    0.000067
	peak: X=120, Y=1782.31

`
	actual = string(stderr)
	if expected != actual {
		t.Errorf("Expected\n%s\nbut was\n%s", expected, actual)
	}
}

func fakeMain(t *testing.T, args ...string) (stdoutData, stderrData []byte) {
	stdout, err := ioutil.TempFile(os.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()

	stderr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	os.Stdout = stdout
	os.Stderr = stderr

	os.Args = append([]string{"usl"}, args...)

	main()

	err = stdout.Sync()
	if err != nil {
		t.Error(err)
	}

	err = stderr.Sync()
	if err != nil {
		t.Error(err)
	}

	stdoutData, err = ioutil.ReadFile(stdout.Name())
	if err != nil {
		t.Fatal(err)
	}

	stderrData, err = ioutil.ReadFile(stderr.Name())
	if err != nil {
		t.Fatal(err)
	}

	return
}