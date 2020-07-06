// USL is a modeler for the Universal Scalability Law, which can be used in system testing and
// capacity planning.
//
// As an example, consider doing load testing and capacity planning for an HTTP server. To use USL,
// we must first gather a set of measurements of the system. These measurements will consist of
// pairs of simultaneous measurements of the independent and dependent variables. With an HTTP
// server, it might be tempting to use the rate as the independent variable, but this is a mistake.
// The rate of requests being handled by the server is actually itself a dependent variable of two
// other independent variables: the number of concurrent users and the rate at which users send
// requests.
//
// As we do our capacity planning, we make the observation that users of our system do ~10 req/sec.
// (Or, more commonly, we assume this based on a hunch.) By holding this constant, we leave the
// number of concurrent users as the single remaining independent variable.
//
// Our load testing, then, should consist of running a series of tests with an increasing number of
// simulated users, each performing ~10 req/sec. While the number of users to test with depends
// heavily on your system, you should be testing at least six different concurrency levels. You
// should do one test with a single user in order to determine the performance of an uncontended
// system.
//
// After our load testing is done, we should have a CSV file which consists of a series of
// (concurrency, throughput) pairs of measurements:
//
//      1,65
//      18,996
//      36,1652
//      72,1853
//      108,1829
//      144,1775
//      216,1702
//
// We can then run the USL binary:
//
//		usl -in data.csv
//
// USL parses the given CSV file as a series of (concurrency, throughput) points, calculates the USL
// parameters using quadratic regression, and then prints out the details of the model:
//
//     URL parameters: σ=0.02772985648395876, κ=0.00010434289088915312, λ=89.98778453648904
//         max throughput: 1883.7622524836281, max concurrency: 96
//         contention constrained
//
// Among the details here we see two things worth noting. First, the system appears to be
// constrained by contention, so optimization work should be focused mostly on removing locks, etc.
// Second, the peak throughput of the system is expected to occur at 96 concurrent users, at which
// point the system will be expected to handle ~1883 req/sec.
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

func main() {
	input := flag.String("in", "", "input file")
	nCol := flag.Int("n_col", 1, "column index of concurrency values")
	rCol := flag.Int("r_col", 2, "column index of latency values")
	skip := flag.Bool("skip_headers", false, "skip the first line")

	log.SetFlags(0) // don't prefix the log statements
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Printf("Usage: usl <-in input.csv> [options] [points...]\n\n")
		flag.PrintDefaults()
	}
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
	log.Printf("URL parameters: σ=%v, κ=%v, λ=%v\n", m.Sigma, m.Kappa, m.Lambda)
	log.Printf("\tmax throughput: %v, max concurrency: %v\n", m.MaxThroughput(), m.MaxConcurrency())

	if m.ContentionConstrained() {
		log.Println("\tcontention constrained")
	}

	if m.CoherencyConstrained() {
		log.Println("\tcoherence constrained")
	}

	if m.Limitless() {
		log.Println("\tlimitless")
	}

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
	measurements := make([]usl.Measurement, 0, 100)

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
		n, x, err := parseLine(i, nCol, rCol, line)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, usl.ConcurrencyAndThroughput(n, x))
	}

	return measurements, nil
}

//nolint:goerr113 // not a package
func parseLine(i, nCol, xCol int, line []string) (float64, float64, error) {
	if len(line) != 2 {
		return 0, 0, fmt.Errorf("invalid line at line %d", i+1)
	}

	n, err := strconv.ParseFloat(line[nCol-1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error at line %d, column %d: %w", i+1, nCol, err)
	}

	x, err := strconv.ParseFloat(line[xCol-1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error at line %d, column %d: %w", i+1, xCol, err)
	}

	return n, x, nil
}
