package usl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConcurrencyAndLatency(t *testing.T) {
	m := ConcurrencyAndLatency(3, 0.6)

	want, got := 3.0, m.Concurrency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.6, m.Latency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 5, m.Throughput
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestConcurrencyAndThroughput(t *testing.T) {
	m := ConcurrencyAndThroughput(3, 5)

	want, got := 3.0, m.Concurrency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.6, m.Latency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 5, m.Throughput
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestThroughputAndLatency(t *testing.T) {
	m := ThroughputAndLatency(5, 0.6)

	want, got := 3.0, m.Concurrency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.6, m.Latency
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 5, m.Throughput
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}
