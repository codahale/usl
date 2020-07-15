package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codahale/usl"
	"github.com/codahale/usl/internal/assert"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParsing(t *testing.T) {
	want := []usl.Measurement{
		usl.ConcurrencyAndThroughput(1, 65),
		usl.ConcurrencyAndThroughput(18, 996),
		usl.ConcurrencyAndThroughput(36, 1652),
		usl.ConcurrencyAndThroughput(72, 1853),
		usl.ConcurrencyAndThroughput(108, 1829),
		usl.ConcurrencyAndThroughput(144, 1775),
		usl.ConcurrencyAndThroughput(216, 1702),
	}

	got, err := parseCSV("example.csv", 1, 2, false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "measurements", want, got,
		cmpopts.EquateApprox(0.001, 0.001))
}

func TestBadLine(t *testing.T) {
	_, _, err := parseLine(0, 1, 2, []string{"funk"})
	if err == nil {
		t.Fatalf("should have failed")
	}
}

func TestBadConcurrency(t *testing.T) {
	_, _, err := parseLine(0, 1, 2, []string{"f", "1"})
	if err == nil {
		t.Fatalf("should have failed")
	}
}

func TestBadThroughput(t *testing.T) {
	_, _, err := parseLine(0, 1, 2, []string{"1", "f"})
	if err == nil {
		t.Fatalf("should have failed")
	}
}

func TestMainRun(t *testing.T) {
	stdout, stderr := fakeMain(t, "-in", "example.csv", "1", "2", "3")

	assert.Equal(t, "stdout",
		`1.000000,89.987785
2.000000,175.083978
3.000000,255.626353
`,
		string(stdout))

	fmt.Println(string(stderr))
	assert.Equal(t, "stderr",
		`URL parameters: σ=0.02772985648395876, κ=0.00010434289088915312, λ=89.98778453648904
	max throughput: 1883.7622524836281, max concurrency: 96
	contention constrained
                                                                          
        |                                                                 
 2.1 k  +                                                                 
 2.0 k  +                   ***X******@***X*********                      
 1.8 k  +              *****                        **X**********         
 1.7 k  +          X **                                                   
 1.6 k  +          **                                                     
 1.5 k  +        **                                                       
 1.3 k  +       *                                                         
 1.2 k  +      *                                                          
 1.1 k  +    X*                                                           
   975  +     *                                  .------------------.     
   853  +    *                                   |****** Predicted  |     
   731  +   *                                    |  X    Actual     |     
   609  +  *                                     |  @    Peak       |     
   487  + *                                      '------------------'     
   366  + *                                                               
   244  +*                                                                
  1122  X----+------+-----+-----+-----+------+-----+-----+-----+-----+-   
            19     38    58    77    96     115   134   154   173   192   

`,
		string(stderr))
}

func fakeMain(t *testing.T, args ...string) ([]byte, []byte) {
	stdout, err := ioutil.TempFile(os.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = stdout.Close() }()

	stderr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = stderr.Close() }()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	os.Stdout = stdout
	os.Stderr = stderr

	os.Args = append([]string{"usl"}, args...)

	main()

	err = stdout.Sync()
	if err != nil {
		t.Error(err)
	}

	err = stderr.Sync()
	if err != nil {
		t.Error(err)
	}

	stdoutData, err := ioutil.ReadFile(stdout.Name())
	if err != nil {
		t.Fatal(err)
	}

	stderrData, err := ioutil.ReadFile(stderr.Name())
	if err != nil {
		t.Fatal(err)
	}

	return stdoutData, stderrData
}
