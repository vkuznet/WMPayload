package main

// Go implementation of WMPayload server
//
// Copyright (c) 2024 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	_ "expvar" // to be used for monitoring, see https://github.com/divan/expvarmon
	"fmt"
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/
	"runtime"
	"time"
)

// version of the code
var gitVersion string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("srv git=%s go=%s date=%s", gitVersion, goVersion, tstamp)
}

func main() {
	server()
}
