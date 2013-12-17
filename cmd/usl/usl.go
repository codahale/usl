// USL is a modeler for the Universal Scalability Law, which can be used in
// system testing and capacity planning.
//
// Usage:
//
//		usl -in data.csv [x ...]
//
// USL parses the given CSV file as a series of (x, y) points, calculates the
// USL parameters using quadratic regression, and then evaluates any given data
// points.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/codahale/usl"
)

var (
	input = flag.String("in", "", "input file")
)

func main() {
	log.SetFlags(0) // don't prefix the log statements
	log.SetOutput(os.Stderr)
	flag.Parse()

	if len(*input) == 0 {
		log.Fatal("No input files provided.")
	}

	measurements, err := parseCSV(*input)
	if err != nil {
		log.Fatal(err)
	}

	m, err := usl.Build(measurements)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Model:")
	log.Printf("\tα:    %f\n", m.Alpha)
	log.Printf("\tβ:    %f\n", m.Beta)
	log.Printf("\tpeak: X=%.0f, Y=%2.2f\n", m.Peak, m.Predict(m.Peak))
	log.Println()

	for _, s := range flag.Args() {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%f,%f\n", x, m.Predict(x))
	}
}

func parseCSV(filename string) (usl.MeasurementSet, error) {
	measurements := make(usl.MeasurementSet, 0)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, line := range lines {
		m, err := parseLine(i, line)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, m)
	}

	return measurements, nil
}

func parseLine(i int, line []string) (m usl.Measurement, err error) {
	if len(line) != 2 {
		err = fmt.Errorf("invalid line at line %d", i+1)
		return
	}

	m.X, err = strconv.ParseFloat(line[0], 64)
	if err != nil {
		err = fmt.Errorf("%v at line %d, column 1", err, i+1)
		return
	}

	m.Y, err = strconv.ParseFloat(line[1], 64)
	if err != nil {
		err = fmt.Errorf("%v at line %d, column 2", err, i+1)
		return
	}

	return
}
