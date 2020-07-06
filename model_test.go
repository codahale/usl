package usl

import (
	"testing"

	"github.com/codahale/usl/internal/assert"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	//nolint:gochecknoglobals // fine in tests
	epsilon = cmpopts.EquateApprox(0.00001, 0.00001)
)

func TestModel_Kappa(t *testing.T) {
	m := build(t)

	assert.Equal(t, "Kappa", 7.690945e-4, m.Kappa, epsilon)
}

func TestModel_Sigma(t *testing.T) {
	m := build(t)

	assert.Equal(t, "Sigma", 0.02671591, m.Sigma, epsilon)
}

func TestModel_Lambda(t *testing.T) {
	m := build(t)

	assert.Equal(t, "Lambda", 995.6486, m.Lambda, epsilon)
}

func TestModel_MaxConcurrency(t *testing.T) {
	m := build(t)

	assert.Equal(t, "MaxConcurrency", 35.0, m.MaxConcurrency(), epsilon)
}

func TestModel_MaxThroughput(t *testing.T) {
	m := build(t)

	assert.Equal(t, "MaxThroughput", 12341.745415132369, m.MaxThroughput(), epsilon)
}

func TestModel_CoherencyConstrained(t *testing.T) {
	m := build(t)

	assert.Equal(t, "CoherencyConstrained", false, m.CoherencyConstrained())
}

func TestModel_ContentionConstrained(t *testing.T) {
	m := build(t)

	assert.Equal(t, "ContentionConstrained", true, m.ContentionConstrained())
}

func TestModel_LatencyAtConcurrency(t *testing.T) {
	m := build(t)

	assert.Equal(t, "R(N=1)", 0.0010043702853760425, m.LatencyAtConcurrency(1), epsilon)
	assert.Equal(t, "R(N=20)", 0.0018077244276309343, m.LatencyAtConcurrency(20), epsilon)
	assert.Equal(t, "R(N=35)", 0.0028359035794958197, m.LatencyAtConcurrency(35), epsilon)
}

func TestModel_ThroughputAtConcurrency(t *testing.T) {
	m := build(t)

	assert.Equal(t, "X(N=1)", 995.648772003358, m.ThroughputAtConcurrency(1), epsilon)
	assert.Equal(t, "X(N=20)", 11063.63312570436, m.ThroughputAtConcurrency(20), epsilon)
	assert.Equal(t, "X(N=35)", 12341.745655201905, m.ThroughputAtConcurrency(35), epsilon)
}

func TestModel_ConcurrencyAtThroughput(t *testing.T) {
	m := build(t)

	assert.Equal(t, "N(X=955)", 0.9580998829620233, m.ConcurrencyAtThroughput(955), epsilon)
	assert.Equal(t, "N(X=11048)", 15.350435172752203, m.ConcurrencyAtThroughput(11048), epsilon)
	assert.Equal(t, "N(X=12201)", 17.73220762025387, m.ConcurrencyAtThroughput(12201), epsilon)
}

func TestModel_ThroughputAtLatency(t *testing.T) {
	m := &Model{Sigma: 0.06, Kappa: 0.06, Lambda: 40}

	assert.Equal(t, "X(R=0.03)", 69.38886664887109, m.ThroughputAtLatency(0.03), epsilon)
	assert.Equal(t, "X(R=0.04", 82.91561975888501, m.ThroughputAtLatency(0.04), epsilon)
	assert.Equal(t, "X(R=0.05)", 84.06346808612327, m.ThroughputAtLatency(0.05), epsilon)
}

func TestModel_LatencyAtThroughput(t *testing.T) {
	m := &Model{Sigma: 0.06, Kappa: 0.06, Lambda: 40}

	assert.Equal(t, "R(N=400)", 0.05875, m.LatencyAtThroughput(400), epsilon)
	assert.Equal(t, "R(N=500", 0.094, m.LatencyAtThroughput(500), epsilon)
	assert.Equal(t, "R(N=600)", 0.235, m.LatencyAtThroughput(600), epsilon)
}

func TestModel_ConcurrencyAtLatency(t *testing.T) {
	m, err := Build(measurements[:10])
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "N(R=0.0012)", 7.230628979597649, m.ConcurrencyAtLatency(0.0012), epsilon)
	assert.Equal(t, "N(R=0.0016)", 20.25106409917121, m.ConcurrencyAtLatency(0.0016), epsilon)
	assert.Equal(t, "N(R=0.0020)", 29.88889360938781, m.ConcurrencyAtLatency(0.0020), epsilon)
}

func TestModel_Limitless(t *testing.T) {
	m := &Model{Sigma: 1, Lambda: 40}

	assert.Equal(t, "Limitless", true, m.Limitless())

	m = build(t)

	assert.Equal(t, "Limitless", false, m.Limitless())
}

func BenchmarkBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		build(b)
	}
}

//nolint:gochecknoglobals // fine in tests
var measurements = []Measurement{
	ConcurrencyAndThroughput(1, 955.16),
	ConcurrencyAndThroughput(2, 1878.91),
	ConcurrencyAndThroughput(3, 2688.01),
	ConcurrencyAndThroughput(4, 3548.68),
	ConcurrencyAndThroughput(5, 4315.54),
	ConcurrencyAndThroughput(6, 5130.43),
	ConcurrencyAndThroughput(7, 5931.37),
	ConcurrencyAndThroughput(8, 6531.08),
	ConcurrencyAndThroughput(9, 7219.8),
	ConcurrencyAndThroughput(10, 7867.61),
	ConcurrencyAndThroughput(11, 8278.71),
	ConcurrencyAndThroughput(12, 8646.7),
	ConcurrencyAndThroughput(13, 9047.84),
	ConcurrencyAndThroughput(14, 9426.55),
	ConcurrencyAndThroughput(15, 9645.37),
	ConcurrencyAndThroughput(16, 9897.24),
	ConcurrencyAndThroughput(17, 10097.6),
	ConcurrencyAndThroughput(18, 10240.5),
	ConcurrencyAndThroughput(19, 10532.39),
	ConcurrencyAndThroughput(20, 10798.52),
	ConcurrencyAndThroughput(21, 11151.43),
	ConcurrencyAndThroughput(22, 11518.63),
	ConcurrencyAndThroughput(23, 11806),
	ConcurrencyAndThroughput(24, 12089.37),
	ConcurrencyAndThroughput(25, 12075.41),
	ConcurrencyAndThroughput(26, 12177.29),
	ConcurrencyAndThroughput(27, 12211.41),
	ConcurrencyAndThroughput(28, 12158.93),
	ConcurrencyAndThroughput(29, 12155.27),
	ConcurrencyAndThroughput(30, 12118.04),
	ConcurrencyAndThroughput(31, 12140.4),
	ConcurrencyAndThroughput(32, 12074.39),
}

func build(t testing.TB) *Model {
	m, err := Build(measurements)
	if err != nil {
		t.Fatal(err)
	}

	return m
}
