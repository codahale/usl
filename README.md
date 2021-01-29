usl
===

`usl` is a Go modeler for [Dr. Neil Gunther][NJG]'s [Universal Scalability Law][USL] as described in
[Baron Schwartz][BS]'s book [Practical Scalability Analysis with the Universal Scalability
Law][PSA].

Given a handful of measurements of any two [Little's Law][LL] parameters--throughput, latency, and
concurrency--the [USL][USL] allows you to make predictions about any of those parameters' values
given an arbitrary value for any another parameter. For example, given a set of measurements of
concurrency and throughput, the [USL][USL] will allow you to predict what a system's average latency
will look like at a particular throughput, or how many servers you'll need to process requests and
stay under your SLA's latency requirements.

The model coefficients and predictions should be within 0.02% of those listed in the book.

## How to use this

As an example, consider doing load testing and capacity planning for an HTTP server. To model the
behavior of the system using the [USL][USL], you must first gather a set of measurements of the
system. These measurements must be of two of the three parameters of [Little's Law][LL]: mean
response time (in seconds), throughput (in requests per second), and concurrency (i.e. the number of
concurrent clients).

Because response time tends to be a property of load (i.e. it rises as throughput or concurrency
rises), the dependent variable in your tests should be mean response time. This leaves either
throughput or concurrency as your independent variable, but thanks to [Little's Law][LL] it doesn't
matter which one you use. For the purposes of discussion, let's say you measure throughput as a
function of the number of concurrent clients working at a fixed rate (e.g. you used
[`wrk2`][wrk2]).

After you're done load testing, you should have a set of measurements shaped like this:

|concurrency|throughput|
|-----------|----------|
|          1|        65|
|         18|       996|
|         36|      1652|
|         72|      1853|
|        108|      1829|
|        144|      1775|
|        216|      1702|

Now you can build a model and begin estimating things.

### As A CLI Tool

```
$ go get github.com/codahale/usl/cmd/usl
```

```
$ cat measurements.csv
1,65
18,996
36,1652
72,1853
etc.
```

```
$ usl measurements.csv 10 50 100 150 200 250 300
USL parameters: σ=0.0277299, κ=0.000104343, λ=89.9878
	max throughput: 1883.76, max concurrency: 96
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

10.000000,714.778987
50.000000,1721.000603
100.000000,1883.278957
150.000000,1808.481681
200.000000,1686.571797
250.000000,1562.279331
300.000000,1447.463805
```

### As A Go Library


```go
import (
	"fmt"

	"github.com/codahale/usl"
)

func main() {
	measurements := []usl.Measurement{
		usl.ConcurrencyAndThroughput(1, 955.16),
		usl.ConcurrencyAndThroughput(2, 1878.91),
		usl.ConcurrencyAndThroughput(3, 2688.01), // etc
	}

	model := usl.Build(measurements)
	for n := 10; n < 200; n += 10 {
		fmt.Printf("At %d concurrent clients, expect %f req/sec\n",
			n, model.ThroughputAtConcurrency(float64(n)))
	}
}
```

## Performance

Building models is pretty fast:

```
pkg: github.com/codahale/usl
BenchmarkBuild-8   	    2242	    500232 ns/op
```

## Further reading

I strongly recommend [Practical Scalability Analysis with the Universal Scalability Law][PSA], a
free e-book by [Baron Schwartz][BS], author of [High Performance MySQL][MySQL] and CEO of
[VividCortex][VC]. Trying to use this library without actually understanding the concepts behind
[Little's Law][LL], [Amdahl's Law][AL], and the [Universal Scalability Law][USL] will be difficult
and potentially misleading.

I also [wrote a blog post about my Java implementation of USL][usl4j].

## License

Copyright © 2021 Coda Hale

Distributed under the Apache License 2.0.

[NJG]: http://www.perfdynamics.com/Bio/njg.html
[AL]: https://en.wikipedia.org/wiki/Amdahl%27s_law
[LL]: https://en.wikipedia.org/wiki/Little%27s_law
[PSA]: https://www.vividcortex.com/resources/universal-scalability-law/
[USL]: http://www.perfdynamics.com/Manifesto/USLscalability.html
[BS]: https://www.xaprb.com/
[MySQL]: http://shop.oreilly.com/product/0636920022343.do
[VC]: https://www.vividcortex.com/
[wrk2]: https://github.com/giltene/wrk2
[usl4j]: https://codahale.com/usl4j-and-you/
