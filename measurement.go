package usl

import (
	"fmt"
	"time"
)

// Measurement is a simultaneous measurement of at least two of the parameters of Little's Law:
// concurrency, throughput, and latency. The third parameter is inferred from the other two.
type Measurement struct {
	Concurrency float64 // The average number of concurrent events.
	Throughput  float64 // The long-term average arrival rate of events, in events/sec.
	Latency     float64 // The average duration of events in seconds.
}

func (m *Measurement) String() string {
	return fmt.Sprintf("(n=%v,x=%v,r=%v)", m.Concurrency, m.Throughput, m.Latency)
}

// ConcurrencyAndLatency returns a measurement of a system's latency at a given level of
// concurrency. The throughput of the system is derived via Little's Law.
func ConcurrencyAndLatency(n uint64, r time.Duration) Measurement {
	return Measurement{
		Concurrency: float64(n),               // L
		Throughput:  float64(n) / r.Seconds(), // λ=L/W
		Latency:     r.Seconds(),              // W
	}
}

// ConcurrencyAndThroughput returns a measurement of a system's throughput at a given level of
// concurrency. The latency of the system is derived via Little's Law.
func ConcurrencyAndThroughput(n uint64, x float64) Measurement {
	return Measurement{
		Concurrency: float64(n),     // L
		Throughput:  x,              // λ
		Latency:     float64(n) / x, // W=L/λ
	}
}

// ThroughputAndLatency returns a measurement of a system's latency at a given level of throughput.
// The concurrency of the system is derived via Little's Law.
func ThroughputAndLatency(x float64, r time.Duration) Measurement {
	return Measurement{
		Concurrency: x * r.Seconds(), // L=λW
		Throughput:  x,               // λ
		Latency:     r.Seconds(),     // W
	}
}
