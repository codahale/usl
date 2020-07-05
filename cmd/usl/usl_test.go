package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/codahale/usl"
)

func TestParsing(t *testing.T) {
	expected := []usl.Measurement{
		{Concurrency: 1, Throughput: 65},
		{Concurrency: 18, Throughput: 996},
		{Concurrency: 36, Throughput: 1652},
		{Concurrency: 72, Throughput: 1853},
		{Concurrency: 108, Throughput: 1829},
		{Concurrency: 144, Throughput: 1775},
		{Concurrency: 216, Throughput: 1702},
	}

	actual, err := parseCSV("example.csv", 1, 2, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("Expected %v measurements, but was %v", len(expected), len(actual))
	}

	for i, a := range actual {
		e := expected[i]

		if a.Concurrency != e.Concurrency || a.Throughput != e.Throughput {
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
	expected := `1.000000,89.987785
2.000000,175.083978
3.000000,255.626353
`
	stdout, stderr := fakeMain(t, "-in", "example.csv", "1", "2", "3")

	actual := string(stdout)
	if expected != actual {
		t.Errorf("Expected\n%s\nbut was\n%s", expected, actual)
	}

	expected = `Model{σ=0.02772985648395876, κ=0.00010434289088915312, λ=89.98778453648904}

`
	actual = string(stderr)
	if expected != actual {
		t.Errorf("Expected\n%q\nbut was\n%q", expected, actual)
	}
}

func fakeMain(t *testing.T, args ...string) (stdoutData, stderrData []byte) {
	stdout, err := ioutil.TempFile(os.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = stdout.Close() }()

	stderr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = stderr.Close() }()

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
