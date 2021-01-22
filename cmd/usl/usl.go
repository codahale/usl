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
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
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

func run() error {
	//nolint:maligned // ordering of fields matters
	var cli struct {
		InputPath         string           `arg:"" type:"existingfile" help:"The CSV file measurements of the system."`
		Predictions       []float64        `arg:"" optional:"" help:"Predict throughput at the given concurrency levels."`
		ConcurrencyColumn int              `short:"N" default:"1" help:"The column index of concurrency values."`
		LatencyColumn     int              `short:"R" default:"2" help:"The column index of latency values."`
		SkipHeaders       bool             `default:"false" help:"Skip the first line of the file."`
		Width             int              `short:"W" default:"74" help:"The width of the graph in chars."`
		Height            int              `short:"H" default:"20" help:"The height of the graph in chars."`
		NoGraph           bool             `default:"false" help:"Don't display the graph.'"`
		Version           kong.VersionFlag `help:"Display the application version."`
	}

	ctx := kong.Parse(&cli, kong.Vars{"version": version})
	if ctx.Error != nil {
		_, _ = fmt.Fprintln(os.Stderr, ctx.Error)
		os.Exit(1)
	}

	measurements, err := parseCSV(cli.InputPath, cli.ConcurrencyColumn, cli.LatencyColumn, cli.SkipHeaders)
	if err != nil {
		return fmt.Errorf("error parsing %q: %w", cli.InputPath, err)
	}

	m, err := usl.Build(measurements)
	if err != nil {
		return err
	}

	printModel(m, measurements, cli.NoGraph, cli.Width, cli.Height)

	printPredictions(m, cli.Predictions)

	return nil
}

func printModel(m *usl.Model, measurements []usl.Measurement, noGraph bool, width, height int) {
	_, _ = fmt.Fprintf(os.Stderr, "USL parameters: σ=%.6g, κ=%.6g, λ=%.6g\n", m.Sigma, m.Kappa, m.Lambda)
	_, _ = fmt.Fprintf(os.Stderr, "\tmax throughput: %.6g, max concurrency: %.6g\n", m.MaxThroughput(), m.MaxConcurrency())

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

func printPredictions(m *usl.Model, args []float64) {
	for _, n := range args {
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

var version = "dev"
