package usl

import (
	"testing"
	"time"

	"github.com/codahale/gubbins/assert"
)

func TestMeasurement_String(t *testing.T) {
	t.Parallel()

	m := Measurement{Concurrency: 1, Throughput: 2, Latency: 3}

	assert.Equal(t, "String", "(n=1,x=2,r=3)", m.String())
}

func TestConcurrencyAndLatency(t *testing.T) {
	t.Parallel()

	m := ConcurrencyAndLatency(3, 600*time.Millisecond)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}

func TestConcurrencyAndThroughput(t *testing.T) {
	t.Parallel()

	m := ConcurrencyAndThroughput(3, 5)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}

func TestThroughputAndLatency(t *testing.T) {
	t.Parallel()

	m := ThroughputAndLatency(5, 600*time.Millisecond)

	assert.Equal(t, "Concurrency", 3.0, m.Concurrency, epsilon)
	assert.Equal(t, "Latency", 0.6, m.Latency, epsilon)
	assert.Equal(t, "Throughput", 5.0, m.Throughput, epsilon)
}
