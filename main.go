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

// tagVersion of the code shows git tag
var tagVersion string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("wmpayload server tag=%s git=%s go=%s date=%s", tagVersion, gitVersion, goVersion, tstamp)
}

func main() {
	server()
}
