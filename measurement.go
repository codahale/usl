package usl

import "fmt"

// Measurement is a simultaneous measurement of at least two of the parameters of Little's Law:
// concurrency, throughput, and latency. The third parameter is inferred from the other two.
type Measurement struct {
	Concurrency float64 // The average number of concurrent events.
	Throughput  float64 // The long-term average arrival rate of events.
	Latency     float64 // The average duration of events.
}

func (m *Measurement) String() string {
	return fmt.Sprintf("(n=%v,x=%v,r=%v)", m.Concurrency, m.Throughput, m.Latency)
}

// ConcurrencyAndLatency returns a measurement of a system's latency at a given level of
// concurrency. The throughput of the system is derived via Little's Law.
func ConcurrencyAndLatency(n, r float64) Measurement {
	return Measurement{
		Concurrency: n,     // L
		Throughput:  n / r, // λ=L/W
		Latency:     r,     // W
	}
}

// ConcurrencyAndThroughput returns a measurement of a system's throughput at a given level of
// concurrency. The latency of the system is derived via Little's Law.
func ConcurrencyAndThroughput(n, x float64) Measurement {
	return Measurement{
		Concurrency: n,     // L
		Throughput:  x,     // λ
		Latency:     n / x, // W=L/λ
	}
}

// ThroughputAndLatency returns a measurement of a system's latency at a given level of throughput.
// The concurrency of the system is derived via Little's Law.
func ThroughputAndLatency(x, r float64) Measurement {
	return Measurement{
		Concurrency: x * r, // L=λW
		Throughput:  x,     // λ
		Latency:     r,     // W
	}
}
