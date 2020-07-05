package usl

// Measurement is a simultaneous measurement of at least two of the parameters of Little's Law:
// concurrency, throughput, and latency. The third parameter is inferred from the other two.
type Measurement struct {
	Concurrency float64
	Throughput  float64
	Latency     float64
}

// ConcurrencyAndLatency returns a measurement of a system's latency at a given level of
// concurrency.
func ConcurrencyAndLatency(n, r float64) Measurement {
	return Measurement{
		Concurrency: n,
		Throughput:  n / r,
		Latency:     r,
	}
}

// ConcurrencyAndThroughput returns a measurement of a system's throughput at a given level of
// concurrency.
func ConcurrencyAndThroughput(n, x float64) Measurement {
	return Measurement{
		Concurrency: n,
		Throughput:  x,
		Latency:     n / x,
	}
}

// ThroughputAndLatency returns a measurement of a system's latency at a given level of throughput.
func ThroughputAndLatency(x, r float64) Measurement {
	return Measurement{
		Concurrency: x * r,
		Throughput:  x,
		Latency:     r,
	}
}
