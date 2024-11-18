### WMPayload benchmarks
The benchmark was done using Go standard library, see this
[document](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
on how to write them.

We provide two shell scripts:
- `benchmark_insert.sh` to inject data into WMPayload server
- and, `benchmark_search.sh` to search for data in WMPayload server.
Both scripts by default use 100 concurrent clients.

To perform the benchmark we need two terminal windows (or tabs). One
where we run the server and another where we'll run a client benchmark script.
We also configured our server to point to local instance of MongoDB with the
following configuration:
```
{
  "use_https": false,
  "port": 7111,
  "cert_file": "",
  "key_file": "",
  "json_file": "wma.json",
  "mongo_uri": "mongodb://localhost:8230",
  "mongo_database": "fasthttp",
  "mongo_collection": "test"
}
```
Using this configuration we start server in plain HTTP mode, and use 
`wma.json` file which contains one of the ReqMgr2 JSON documents. The content
of the document is irrelevant as we only add to it uuid and re-use the same
content for all injections and search queries.

Please refer to `benchmark_test.go` for actual benchmark implementation.

Here we provide details of performed steps along with obtained results:


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
BenchmarkInsert-8   	    1419	    840000 ns/op	         1.000 success_rate	   60121 B/op	     666 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	1.463s
```

The insert takes 0.5 ms per HTTP call for simple JSON (see `benchmark_test.go`)
and 0.8ms using ReqMgr2
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
BenchmarkSearch-8   	      14	  75389929 ns/op	         1.000 success_rate	102658732 B/op	     238 allocs/op
PASS
ok  	github.com/vkuznet/wmpayload	4.003s
```
So, it is 75ms, and 100MB per operation.

We summarize our results in the following table (please note, all tests done
using JSON data-format):
| operation | document | req/sec | bytes/operation | memory allocations |
|-----------|----------|---------|-----------------|--------------------|
| write     | auto-gen | 0.5ms   | 12KB  | 197 |
| write     | ReqMgr2  | 0.8ms   | 60KB  | 666 |
| read      | auto-gen | 0.2ms   | 12KB  | 124 |
| read      | ReqMgr2  | 75ms    | 102MB | 238 |

Please note
