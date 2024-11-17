### Benchmarks

```
# start the server
./wmpayload -config config.json
2024/11/17 16:22:35 Starting HTTP server on port 7111
```

```
# run insert benchmark

goos: darwin
goarch: arm64
pkg: github.com/dmwm/wmpayload
cpu: Apple M2
BenchmarkInsert-8   	    1840	    598964 ns/op	         1.000 success_rate	   22975 B/op	     197 allocs/op
PASS
ok  	github.com/dmwm/wmpayload	1.445s
```

The insert takes 0.5 ms per HTTP call with 22KB per operation.

```
# run search benchmark
./benchmark_search.sh

goos: darwin
goarch: arm64
pkg: github.com/dmwm/wmpayload
cpu: Apple M2
BenchmarkSearch-8   	     123	   8854730 ns/op	         1.000 success_rate	 8414498 B/op	     177 allocs/op
PASS
ok  	github.com/dmwm/wmpayload	2.561s

```
The search takes 8 ms per query, this rate is poor due to missing index. Let's insert
new index in MongoDB:

```
mongo --port xxxxx
db.test.createIndex({ id: 1 })
```
and, now we can re-run our search benchmark test:
```
# run search benchmark
./benchmark_search.sh

goos: darwin
goarch: arm64
pkg: github.com/dmwm/wmpayload
cpu: Apple M2
BenchmarkSearch-8   	    7771	    249416 ns/op	         0 success_rate	   11415 B/op	     124 allocs/op
PASS
ok  	github.com/dmwm/wmpayload	2.728s
```
Now, the search takes 0.2 ms (10+ times better).
