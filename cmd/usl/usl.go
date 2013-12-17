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

	measurements, err := parse(*input)
	if err != nil {
		log.Fatal(err)
	}

	model, err := usl.Build(measurements)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Model:")
	log.Printf("\tα     = %f\n", model.Alpha)
	log.Printf("\tβ     = %f\n", model.Beta)
	log.Printf("\tN max = %d\n", model.Nmax)
	log.Println()

	for _, s := range flag.Args() {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%f,%f\n", x, model.Predict(x))
	}
}

func parse(filename string) (usl.MeasurementSet, error) {
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
		if len(line) != 2 {
			return nil, fmt.Errorf("invalid line at line %d", i)
		}

		x, err := strconv.ParseFloat(line[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%v at line %d, column 0", err, i)
		}

		y, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%v at line %d, column 1", err, i)
		}

		measurements = append(measurements, usl.Measurement{X: x, Y: y})
	}

	return measurements, nil
}
