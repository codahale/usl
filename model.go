package main

import (
	"fmt"
	"math"

	"code.google.com/p/gomatrix/matrix"
)

type model struct {
	alpha, beta, y float64
	nmax, nopt     int
}

func (m model) String() string {
	return fmt.Sprintf("α=%f/β=%f/Nmax=%d/Nopt=%d", m.alpha, m.beta, m.nmax, m.nopt)
}

func (m model) predict(n int) float64 {
	x := float64(n)
	c := x / (1 + (m.alpha * (x - 1)) + (m.beta * x * (x - 1)))
	return c * m.y
}

func analyze(points map[int]float64) (m model, err error) {
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

	m.alpha = math.Abs(c[2] - c[1])
	m.beta = math.Abs(c[2])
	m.y = y1
	m.nmax = int(math.Floor(math.Sqrt((1 - m.alpha) / m.beta)))
	m.nopt = int(math.Ceil(1 / m.alpha))

	return
}
