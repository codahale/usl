// Package usl provides functionality to build Universal Scalability Law models
// from sets of observed measurements.
package usl

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"code.google.com/p/gomatrix/matrix"
)

var (
	// ErrInsufficientMeasurements is returned when less than 6 measurements
	// were provided.
	ErrInsufficientMeasurements = errors.New("usl: need at least 6 measurements")
)

// Measurement is a simultaneous measurement of both an independent variable and
// a dependent variable.
type Measurement struct {
	X float64 // X is a measurement of the independent variable.
	Y float64 // Y is a simultaneous measurement of the dependent variable.
}

// MeasurementSet is a sortable set of measurements.
type MeasurementSet []Measurement

// Len is the number of measurements in the set.
func (m MeasurementSet) Len() int {
	return len(m)
}

// Less reports whether the measurement with index i should sort before the
// measurement with index j.
func (m MeasurementSet) Less(i, j int) bool {
	return m[i].X < m[j].X
}

// Swap swaps the measurements with indexes i and j.
func (m MeasurementSet) Swap(i, j int) {
	x := m[i]
	m[i] = m[j]
	m[j] = x
}

// Model is a Universal Scalability Law model.
type Model struct {
	Alpha float64 // Alpha represents the levels of contention.
	Beta  float64 // Beta represents the coherency delay.
	Y     float64 // Y is the system's unloaded behavior.
	Peak  float64 // Peak is the usage level at which output is maximized.
}

// Predict returns the predicted value at a given utilization level.
func (m Model) Predict(x float64) float64 {
	c := x / (1 + (m.Alpha * (x - 1)) + (m.Beta * x * (x - 1)))
	return c * m.Y
}

func (m Model) String() string {
	var a, b string
	if m.Alpha > m.Beta {
		a = " (constrained by contention effects)"
	} else if m.Alpha < m.Beta {
		b = " (constrained by coherency effects)"
	}

	return fmt.Sprintf(
		"Model:\n\tα:    %f%s\n\tβ:    %f%s\n\tpeak: X=%.0f, Y=%2.2f",
		m.Alpha, a, m.Beta, b, m.Peak, m.Predict(m.Peak),
	)
}

// Build returns a model whose parameters are generated from the given
// measurements.
func Build(measurements MeasurementSet) (m Model, err error) {
	if len(measurements) < 6 {
		err = ErrInsufficientMeasurements
		return
	}

	sort.Sort(measurements)
	y0 := measurements[0].Y / measurements[0].X

	xs := make([]float64, 0, len(measurements))
	ys := make([]float64, 0, len(measurements))

	for _, m := range measurements {
		xs = append(xs, m.X-1)
		ys = append(ys, (m.X/(m.Y/y0))-1)
	}

	var c [3]float64 // do quadratic regression

	y := matrix.MakeDenseMatrix(ys, len(xs), 1)
	x := matrix.Zeros(len(ys), len(c))
	for i := 0; i < len(xs); i++ {
		ip := 1.0
		for j := 0; j < len(c); j++ {
			x.Set(i, j, ip)
			ip *= xs[i]
		}
	}

	q, r := x.QR()
	qty, err := q.Transpose().Times(y)
	if err != nil {
		return
	}

	for i := len(c) - 1; i >= 0; i-- {
		c[i] = qty.Get(i, 0)
		for j := i + 1; j < len(c); j++ {
			c[i] -= c[j] * r.Get(i, j)
		}
		c[i] /= r.Get(i, i)
	}

	m.Alpha = math.Abs(c[2] - c[1])
	m.Beta = math.Abs(c[2])
	m.Y = y0
	m.Peak = math.Floor(math.Sqrt((1 - m.Alpha) / m.Beta))

	return
}
