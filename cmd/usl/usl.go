// USL is a modeler for the Universal Scalability Law, which can be used in
// system testing and capacity planning.
//
// As an example, consider doing load testing and capacity planning for an HTTP
// server. To use USL, we must first gather a set of measurements of the system.
// These measurements will consist of pairs of simultaneous measurements of the
// independent and dependent variables. With an HTTP server, it might be tempting
// to use the rate as the independent variable, but this is a mistake. The rate
// of requests being handled by the server is actually itself a dependent
// variable of two other independent variables: the number of concurrent users
// and the rate at which users send requests.
//
// As we do our capacity planning, we make the observation that users of our
// system do ~10 req/sec. (Or, more commonly, we assume this based on a hunch.)
// By holding this constant, we leave the number of concurrent users as the
// single remaining independent variable.
//
// Our load testing, then, should consist of running a series of tests with an
// increasing number of simulated users, each performing ~10 req/sec. While the
// number of users to test with depends heavily on your system, you should be
// testing at least six different concurrency levels. You should do one test with
// a single user in order to determine the performance of an uncontended system.
//
// After our load testing is done, we should have a CSV file which consists of
// a series of (concurrency, throughput) pairs of measurements:
//
//		1,4227
//		2,8382
//		4,16479
//		8,31856
//		16,59564
//		32,104462
//		64,162985
//
// We can then run the USL binary:
//
//		usl -in data.csv
//
// USL parses the given CSV file as a series of (concurrency, throughput) points, calculates the
// USL parameters using quadratic regression, and then prints out the details of
// the model:
//
//		Model:
//				α:    0.008550 (constrained by contention effects)
//				β:    0.000030
//				peak: Concurrency=181, Y=217458.30
//
// Among the details here we see two things worth noting. First, the system
// appears to be constrained by contention, so optimization work should be
// focused mostly on removing locks, etc. Second, the peak throughput of the
// system is expected to occur at 181 concurrent users, at which point the system
// will be expected to handle ~217K req/sec.
//
// (These numbers are made up, so don't sweat them.)
//
// Finally, we can provide USL a series of additional data points to provide
// estimates for:
//
//		usl -in data.csv 128 256 512
//
// USL will output the data in CSV format on STDOUT.
//
// For more information, see http://www.perfdynamics.com/Manifesto/USLscalability.html.
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
	nCol  = flag.Int("n_col", 1, "column index of concurrency values")
	rCol  = flag.Int("r_col", 2, "column index of latency values")
	skip  = flag.Bool("skip_headers", false, "skip the first line")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: usl <-in input.csv> [options] [points...]\n\n")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0) // don't prefix the log statements
	log.SetOutput(os.Stderr)
	flag.Parse()

	if len(*input) == 0 {
		log.Fatal("No input files provided.")
	}

	measurements, err := parseCSV(*input, *nCol, *rCol, *skip)
	if err != nil {
		log.Fatal(err)
	}

	m, err := usl.Build(measurements)
	if err != nil {
		log.Fatal(err)
	}

	printModel(m)

	printPredictions(m)
}

func printModel(m *usl.Model) {
	log.Println(m)
	log.Println()
}

func printPredictions(m *usl.Model) {
	for _, s := range flag.Args() {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%f,%f\n", n, m.ThroughputAtConcurrency(n))
	}
}

func parseCSV(filename string, nCol, rCol int, skipHeaders bool) ([]usl.Measurement, error) {
	var measurements []usl.Measurement

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if skipHeaders {
		lines = lines[1:]
	}

	for i, line := range lines {
		m, err := parseLine(i, nCol, rCol, line)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, m)
	}

	return measurements, nil
}

func parseLine(i, xCol, yCol int, line []string) (m usl.Measurement, err error) {
	if len(line) != 2 {
		err = fmt.Errorf("invalid line at line %d", i+1)
		return
	}

	m.Concurrency, err = strconv.ParseFloat(line[xCol-1], 64)
	if err != nil {
		err = fmt.Errorf("%v at line %d, column %d", err, i+1, xCol)
		return
	}

	m.Throughput, err = strconv.ParseFloat(line[yCol-1], 64)
	if err != nil {
		err = fmt.Errorf("%v at line %d, column %d", err, i+1, yCol)
		return
	}

	return
}
