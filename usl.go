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

	points, err := parse(*input)
	if err != nil {
		log.Fatal(err)
	}

	model, err := analyze(points)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Model:")
	log.Printf("\tα     = %f\n", model.alpha)
	log.Printf("\tβ     = %f\n", model.alpha)
	log.Printf("\tN max = %d\n", model.nmax)
	log.Printf("\tN opt = %d\n", model.nopt)
	log.Println()

	for _, s := range flag.Args() {
		x, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d,%f\n", x, model.predict(x))
	}
}

func parse(filename string) (map[int]float64, error) {
	points := make(map[int]float64)

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

		x, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("%v at line %d, column 0", err, i)
		}

		y, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%v at line %d, column 1", err, i)
		}

		points[x] = y
	}

	return points, nil
}
