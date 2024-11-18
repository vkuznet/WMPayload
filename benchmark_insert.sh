#!/bin/bash
# Usage: ./insert_benchmark.sh <concurrency>
CONCURRENCY=${1:-100}
# echo "### benchmark db insert"
# echo
# go test -bench=BenchmarkDBInsert -benchmem -parallel $CONCURRENCY -count=1
echo
echo "### benchmark HTTP insert with $CONCURRENCY clients"
echo
go test -bench=BenchmarkInsert -benchmem -parallel $CONCURRENCY -count=1
