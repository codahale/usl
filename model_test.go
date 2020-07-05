package usl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	epsilon = cmpopts.EquateApprox(0.00001, 0.00001)
)

func TestModel_Kappa(t *testing.T) {
	m := build(t)
	want, got := 7.690945e-4, m.Kappa

	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_Sigma(t *testing.T) {
	m := build(t)
	want, got := 0.02671591, m.Sigma

	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_Lambda(t *testing.T) {
	m := build(t)
	want, got := 995.6486, m.Lambda

	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_MaxConcurrency(t *testing.T) {
	m := build(t)
	want, got := 35.0, m.MaxConcurrency()

	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_MaxThroughput(t *testing.T) {
	m := build(t)
	want, got := 12341.745415132369, m.MaxThroughput()

	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_CoherencyConstrained(t *testing.T) {
	m := build(t)
	got := m.CoherencyConstrained()

	if diff := cmp.Diff(false, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_ContentionConstrained(t *testing.T) {
	m := build(t)

	got := m.ContentionConstrained()
	if diff := cmp.Diff(true, got, epsilon); diff != "" {
		t.Fatal(diff)
	}
}

func TestModel_LatencyAtConcurrency(t *testing.T) {
	m := build(t)

	want, got := 0.0010043702853760425, m.LatencyAtConcurrency(1)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.0018077244276309343, m.LatencyAtConcurrency(20)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.0028359035794958197, m.LatencyAtConcurrency(35)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_ThroughputAtConcurrency(t *testing.T) {
	m := build(t)

	want, got := 995.648772003358, m.ThroughputAtConcurrency(1)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 11063.63312570436, m.ThroughputAtConcurrency(20)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 12341.745655201905, m.ThroughputAtConcurrency(35)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_ConcurrencyAtThroughput(t *testing.T) {
	m := build(t)

	want, got := 0.9580998829620233, m.ConcurrencyAtThroughput(955)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 15.350435172752203, m.ConcurrencyAtThroughput(11048)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 17.73220762025387, m.ConcurrencyAtThroughput(12201)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_ThroughputAtLatency(t *testing.T) {
	m := &Model{Sigma: 0.06, Kappa: 0.06, Lambda: 40}

	want, got := 69.38886664887109, m.ThroughputAtLatency(0.03)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 82.91561975888501, m.ThroughputAtLatency(0.04)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 84.06346808612327, m.ThroughputAtLatency(0.05)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_LatencyAtThroughput(t *testing.T) {
	m := &Model{Sigma: 0.06, Kappa: 0.06, Lambda: 40}

	want, got := 0.05875, m.LatencyAtThroughput(400)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.094, m.LatencyAtThroughput(500)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 0.235, m.LatencyAtThroughput(600)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_ConcurrencyAtLatency(t *testing.T) {
	m, err := Build(measurements[:10])
	if err != nil {
		t.Fatal(err)
	}

	want, got := 7.230628979597649, m.ConcurrencyAtLatency(0.0012)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 20.25106409917121, m.ConcurrencyAtLatency(0.0016)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}

	want, got = 29.88889360938781, m.ConcurrencyAtLatency(0.0020)
	if diff := cmp.Diff(want, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func TestModel_Limitless(t *testing.T) {
	m := &Model{Sigma: 1, Lambda: 40}

	got := m.Limitless()
	if diff := cmp.Diff(true, got, epsilon); diff != "" {
		t.Error(diff)
	}

	m = build(t)

	got = m.Limitless()
	if diff := cmp.Diff(false, got, epsilon); diff != "" {
		t.Error(diff)
	}
}

func BenchmarkBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		build(b)
	}
}

var measurements = []Measurement{
	{Concurrency: 1, Throughput: 955.16},
	{Concurrency: 2, Throughput: 1878.91},
	{Concurrency: 3, Throughput: 2688.01},
	{Concurrency: 4, Throughput: 3548.68},
	{Concurrency: 5, Throughput: 4315.54},
	{Concurrency: 6, Throughput: 5130.43},
	{Concurrency: 7, Throughput: 5931.37},
	{Concurrency: 8, Throughput: 6531.08},
	{Concurrency: 9, Throughput: 7219.8},
	{Concurrency: 10, Throughput: 7867.61},
	{Concurrency: 11, Throughput: 8278.71},
	{Concurrency: 12, Throughput: 8646.7},
	{Concurrency: 13, Throughput: 9047.84},
	{Concurrency: 14, Throughput: 9426.55},
	{Concurrency: 15, Throughput: 9645.37},
	{Concurrency: 16, Throughput: 9897.24},
	{Concurrency: 17, Throughput: 10097.6},
	{Concurrency: 18, Throughput: 10240.5},
	{Concurrency: 19, Throughput: 10532.39},
	{Concurrency: 20, Throughput: 10798.52},
	{Concurrency: 21, Throughput: 11151.43},
	{Concurrency: 22, Throughput: 11518.63},
	{Concurrency: 23, Throughput: 11806},
	{Concurrency: 24, Throughput: 12089.37},
	{Concurrency: 25, Throughput: 12075.41},
	{Concurrency: 26, Throughput: 12177.29},
	{Concurrency: 27, Throughput: 12211.41},
	{Concurrency: 28, Throughput: 12158.93},
	{Concurrency: 29, Throughput: 12155.27},
	{Concurrency: 30, Throughput: 12118.04},
	{Concurrency: 31, Throughput: 12140.4},
	{Concurrency: 32, Throughput: 12074.39},
}

func build(t testing.TB) *Model {
	m, err := Build(measurements)
	if err != nil {
		t.Fatal(err)
	}

	return m
}
