### Benchmarks

```
# start the server
./wmpayload -config config.json
2024/11/17 16:22:35 Starting HTTP server on port 7111
```

```
# open up another terminal (tab) and
# run insert benchmark using flat json

goos: darwin
goarch: arm64
pkg: github.com/vkuznet/wmpayload
cpu: Apple M2
BenchmarkInsert-8   	    1840	    598964 ns/op	         1.000 success_rate	   22975 B/op	     197 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	1.445s

# run insert benchmark using ReqMgr2 JSON
goos: darwin
goarch: arm64
pkg: github.com/vkuznet/wmpayload
cpu: Apple M2
BenchmarkInsert-8   	    1088	   1154884 ns/op	         1.000 success_rate	  168339 B/op	    1490 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	1.623s
```

The insert takes 0.5 ms per HTTP call for simple JSON (see `benchmark_test.go`)
and 1.1ms using ReqMgr2
JSON document. The JSON document was take from the following
[URL](https://cmsweb-testbed.cern.ch/reqmgr2/fetch?rid=request-apiccine_SC_PREMIX_GFALStageoutTest_v5_241008_120259_4210) and its size is 12KB.

```
# run search benchmark
./benchmark_search.sh

goos: darwin
goarch: arm64
pkg: github.com/vkuznet/wmpayload
cpu: Apple M2
BenchmarkSearch-8   	     123	   8854730 ns/op	         1.000 success_rate	 8414498 B/op	     177 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	2.561s

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
pkg: github.com/vkuznet/wmpayload
cpu: Apple M2
BenchmarkSearch-8   	    7771	    249416 ns/op	         0 success_rate	   11415 B/op	     124 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	2.728s
```
Now, the search takes 0.2 ms (10+ times better).

And, using ReqMgr2 JSON docs in MongoDB we got the following performance:
```
goos: darwin
goarch: arm64
pkg: github.com/vkuznet/wmpayload
cpu: Apple M2
BenchmarkSearch-8   	       4	 268701615 ns/op	         1.000 success_rate	314660348 B/op	     362 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	4.003s
```
So, it is 268ms, and 300MB per operation.
