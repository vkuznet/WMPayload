#!/bin/bash
# Usage: ./search_benchmark.sh <concurrency>
CONCURRENCY=${1:-100}
echo "### benchmark HTTP search"
echo
go test -bench=BenchmarkSearch -benchmem -parallel $CONCURRENCY -count=1
