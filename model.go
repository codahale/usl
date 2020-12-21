// Package usl provides functionality to build Universal Scalability Law models
// from sets of observed measurements.
package usl

import (
	"fmt"
	"math"

	"github.com/maorshutman/lm"
)

// Model is a Universal Scalability Law model.
type Model struct {
	Sigma  float64 // The model's coefficient of contention, σ.
	Kappa  float64 // The model's coefficient of crosstalk/coherency, κ.
	Lambda float64 // The model's coefficient of performance, λ.
}

func (m *Model) String() string {
	return fmt.Sprintf("Model{σ=%v,κ=%v,λ=%v}", m.Sigma, m.Kappa, m.Lambda)
}

// ThroughputAtConcurrency returns the expected throughput given a number of concurrent events,
// X(N).
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 3".
func (m *Model) ThroughputAtConcurrency(n float64) float64 {
	return (m.Lambda * n) / (1 + (m.Sigma * (n - 1)) + (m.Kappa * n * (n - 1)))
}

// LatencyAtConcurrency returns the expected mean latency given a number of concurrent events,
// R(N).
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 6".
func (m *Model) LatencyAtConcurrency(n float64) float64 {
	return (1 + (m.Sigma * (n - 1)) + (m.Kappa * n * (n - 1))) / m.Lambda
}

// MaxConcurrency returns the maximum expected number of concurrent events the system can handle,
// Nmax.
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 4".
func (m *Model) MaxConcurrency() float64 {
	return math.Floor(math.Sqrt((1 - m.Sigma) / m.Kappa))
}

// MaxThroughput returns the maximum expected throughput the system can handle, Xmax.
func (m Model) MaxThroughput() float64 {
	return m.ThroughputAtConcurrency(m.MaxConcurrency())
}

// LatencyAtThroughput returns the expected mean latency given a throughput, R(X).
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 8".
func (m *Model) LatencyAtThroughput(x float64) float64 {
	return (m.Sigma - 1) / (m.Sigma*x - m.Lambda)
}

// ThroughputAtLatency returns the expected throughput given a mean latency, X(R).
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 9".
func (m *Model) ThroughputAtLatency(r float64) float64 {
	return (math.Sqrt(math.Pow(m.Sigma, 2)+math.Pow(m.Kappa, 2)+
		2*m.Kappa*(2*m.Lambda*r+m.Sigma-2)) - m.Kappa + m.Sigma) / (2.0 * m.Kappa * r)
}

// ConcurrencyAtLatency returns the expected number of concurrent events at a particular mean
// latency, N(R).
//
// See "Practical Scalability Analysis with the Universal Scalability Law, Equation 10".
func (m *Model) ConcurrencyAtLatency(r float64) float64 {
	return (m.Kappa - m.Sigma +
		math.Sqrt(math.Pow(m.Sigma, 2)+
			math.Pow(m.Kappa, 2)+
			2*m.Kappa*((2*m.Lambda*r)+m.Sigma-2))) / (2 * m.Kappa)
}

// ConcurrencyAtThroughput returns the expected number of concurrent events at a particular
// throughput, N(X).
func (m *Model) ConcurrencyAtThroughput(x float64) float64 {
	return m.LatencyAtThroughput(x) * x
}

// ContentionConstrained returns true if the system is constrained by contention.
func (m *Model) ContentionConstrained() bool {
	return m.Sigma > m.Kappa
}

// CoherencyConstrained returns true if the system is constrained by coherency costs.
func (m *Model) CoherencyConstrained() bool {
	return m.Sigma < m.Kappa
}

// Limitless returns true if the system is linearly scalable.
func (m *Model) Limitless() bool {
	return m.Kappa == 0
}

// Build returns a model whose parameters are generated from the given measurements.
//
// Finds a set of coefficients for the equation y = λx/(1+σ(x-1)+κx(x-1)) which best fit the
// observed values using unconstrained least-squares regression. The resulting values for λ, κ, and
// σ are the parameters of the returned model.
func Build(measurements []Measurement) (m *Model, err error) {
	if len(measurements) < minMeasurements {
		return nil, ErrInsufficientMeasurements
	}

	// Calculate an initial guess at the model parameters.
	init := []float64{0.1, 0.01, 0}

	// Use max(x/n) as initial lambda.
	for _, m := range measurements {
		v := m.Throughput / m.Concurrency
		if v > init[2] {
			init[2] = v
		}
	}

	// Calculate the residuals of a possible model.
	f := func(dst, x []float64) {
		model := Model{Sigma: x[0], Kappa: x[1], Lambda: x[2]}

		for i, v := range measurements {
			dst[i] = v.Throughput - model.ThroughputAtConcurrency(v.Concurrency)
		}
	}

	// Formulate an LM problem.
	p := lm.LMProblem{
		Dim:        3,                      // Three parameters in the model.
		Size:       len(measurements),      // Use all measurements to calculate residuals.
		Func:       f,                      // Reduce the residuals of model predictions to observations.
		Jac:        lm.NumJac{Func: f}.Jac, // Approximate the Jacobian by finite differences.
		InitParams: init,                   // Use our initial guesses at parameters.
		Tau:        1e-6,                   // Need a non-zero initial damping factor.
		Eps1:       1e-8,                   // Small but non-zero values here prevent singular matrices.
		Eps2:       1e-8,
	}

	// Calculate the model parameters.
	results, err := lm.LM(p, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build model: %w", err)
	}

	// Return the model.
	return &Model{
		Sigma:  results.X[0],
		Kappa:  results.X[1],
		Lambda: results.X[2],
	}, nil
}

const (
	// minMeasurement is the smallest number of measurements from which a useful model can be
	// created.
	minMeasurements = 6
)

// ErrInsufficientMeasurements is returned when fewer than 6 measurements were provided.
var ErrInsufficientMeasurements = fmt.Errorf("usl: need at least %d measurements", minMeasurements)
