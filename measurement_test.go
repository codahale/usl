package usl

import (
	"testing"

	"github.com/codahale/usl/internal/assert"
)

func TestMeasurement_String(t *testing.T) {
	m := Measurement{Concurrency: 1, Throughput: 2, Latency: 3}

	assert.Equal(t, "String", "(n=1,x=2,r=3)", m.String())
}

func TestConcurrencyAndLatency(t *testing.T) {
	m := ConcurrencyAndLatency(3, 0.6)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}

func TestConcurrencyAndThroughput(t *testing.T) {
	m := ConcurrencyAndThroughput(3, 5)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}

func TestThroughputAndLatency(t *testing.T) {
	m := ThroughputAndLatency(5, 0.6)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}
