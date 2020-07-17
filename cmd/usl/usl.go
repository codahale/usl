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
//     1,65
//     18,996
//     36,1652
//     72,1853
//     108,1829
//     144,1775
//     216,1702
//
// We can then run the USL binary:
//
//     usl data.csv
//
// USL parses the given CSV file as a series of (concurrency, throughput) points, calculates the USL
// parameters using quadratic regression, and then prints out the details of the model, along with a
// graph of the model's predictions and the given measurements.
//
// Finally, we can provide USL a series of additional data points to provide
// estimates for:
//
//     usl data.csv 128 256 512
//
// USL will output the data in CSV format on STDOUT.
//
// For more information, see http://www.perfdynamics.com/Manifesto/USLscalability.html.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/codahale/usl"
	"github.com/vdobler/chart"
	"github.com/vdobler/chart/txtg"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(-1)
	}
}

//nolint:goerr113 // not a package
func run() error {
	nCol := flag.Int("n_col", 1, "column index of concurrency values")
	rCol := flag.Int("r_col", 2, "column index of latency values")
	skipHeaders := flag.Bool("skip_headers", false, "skip the first line")
	width := flag.Int("width", 74, "width of graph")
	height := flag.Int("height", 20, "height of graph")
	noGraph := flag.Bool("no_graph", false, "don't print the graph")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: usl <input.csv> [options] [points...]\n\n")

		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		return fmt.Errorf("no input file provided")
	}

	measurements, err := parseCSV(flag.Arg(0), *nCol, *rCol, *skipHeaders)
	if err != nil {
		return fmt.Errorf("error parsing %w", err)
	}

	m, err := usl.Build(measurements)
	if err != nil {
		return err
	}

	printModel(m, measurements, *noGraph, *width, *height)

	return printPredictions(m, flag.Args()[1:])
}

func printModel(m *usl.Model, measurements []usl.Measurement, noGraph bool, width, height int) {
	_, _ = fmt.Fprintf(os.Stderr, "URL parameters: σ=%v, κ=%v, λ=%v\n", m.Sigma, m.Kappa, m.Lambda)
	_, _ = fmt.Fprintf(os.Stderr, "\tmax throughput: %v, max concurrency: %v\n", m.MaxThroughput(), m.MaxConcurrency())

	if m.ContentionConstrained() {
		_, _ = fmt.Fprintln(os.Stderr, "\tcontention constrained")
	}

	if m.CoherencyConstrained() {
		_, _ = fmt.Fprintln(os.Stderr, "\tcoherence constrained")
	}

	if m.Limitless() {
		_, _ = fmt.Fprintln(os.Stderr, "\tlimitless")
	}

	if !noGraph {
		x := make([]float64, len(measurements))
		y := make([]float64, len(measurements))

		for i, m := range measurements {
			x[i] = m.Concurrency
			y[i] = m.Throughput
		}

		c := chart.ScatterChart{}
		c.Key.Pos = "ibr"
		c.XRange.Fixed(1, m.MaxConcurrency()*2, (m.MaxConcurrency()*2)/10)
		c.YRange.Fixed(0, m.MaxThroughput()*1.1, 0)
		c.NSamples = len(measurements)
		c.AddFunc("Predicted", m.ThroughputAtConcurrency,
			chart.PlotStyleLines, chart.AutoStyle(6, false))
		c.AddDataPair("Actual", x, y, chart.PlotStylePoints, chart.AutoStyle(5, false))
		c.AddDataPair("Peak", []float64{m.MaxConcurrency()}, []float64{m.MaxThroughput()},
			chart.PlotStylePoints, chart.AutoStyle(7, false))

		txt := txtg.New(width, height)
		c.Plot(txt)

		_, _ = fmt.Fprint(os.Stderr, txt)
	}

	_, _ = fmt.Fprintln(os.Stderr)
}

func printPredictions(m *usl.Model, args []string) error {
	for _, s := range args {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		fmt.Printf("%f,%f\n", n, m.ThroughputAtConcurrency(n))
	}

	return nil
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
func parseLine(i, nCol, xCol int, line []string) (uint64, float64, error) {
	if len(line) != 2 {
		return 0, 0, fmt.Errorf("invalid line at line %d", i+1)
	}

	n, err := strconv.ParseUint(line[nCol-1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error at line %d, column %d: %w", i+1, nCol, err)
	}

	x, err := strconv.ParseFloat(line[xCol-1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error at line %d, column %d: %w", i+1, xCol, err)
	}

	return n, x, nil
}
