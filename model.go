package usl

import (
	"fmt"
	"math"

	"code.google.com/p/gomatrix/matrix"
)

type Model struct {
	Alpha, Beta, Y float64
	Nmax, Nopt     int
}

func (m Model) Predict(n int) float64 {
	x := float64(n)
	c := x / (1 + (m.Alpha * (x - 1)) + (m.Beta * x * (x - 1)))
	return c * m.Y
}

// performs quadratic regression on the data points and returns a model
func Analyze(points map[int]float64) (m Model, err error) {
	xs := make([]float64, 0, len(points))
	ys := make([]float64, 0, len(points))

	y1, ok := points[1]
	if !ok {
		err = fmt.Errorf("No point for 1 user")
		return
	}

	for n, y := range points {
		x := float64(n)
		linX := x - 1
		linY := (x / (y / y1)) - 1
		xs = append(xs, linX)
		ys = append(ys, linY)
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
	m.Y = y1
	m.Nmax = int(math.Floor(math.Sqrt((1 - m.Alpha) / m.Beta)))
	m.Nopt = int(math.Ceil(1 / m.Alpha))

	return
}
