GITVERSION=`git rev-parse --short HEAD`
VERSION=`git describe --tags`
TAG := $(shell git tag | sed -e "s,v,," | sort -r | head -n 1)
OS := $(shell uname)
ifeq ($(OS),Darwin)
flags=-ldflags="-s -w -X main.gitVersion=${GITVERSION} -X main.tagVersion=${VERSION}"
else
flags=-ldflags="-s -w -X main.gitVersion=${GITVERSION} -X main.tagVersion=${VERSION} -extldflags -static"
endif


all: build

gorelease:
	goreleaser release --snapshot --clean

build:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg; go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_all: build_darwin_amd64 build_darwin_arm64 build_amd64 build_arm64 build_power8 build_windows_amd64 build_windows_arm64 changes

build_darwin_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_darwin; GOOS=darwin go build -o srv ${flags}
	mv srv srv_darwin_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_darwin_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_darwin; GOARCH=arm64 GOOS=darwin go build -o srv ${flags}
	mv srv srv_darwin_arm64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_linux; GOOS=linux go build -o srv ${flags}
	mv srv srv_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_power8:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_power8; GOARCH=ppc64le GOOS=linux go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv srv_power8

build_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_arm64; GOARCH=arm64 GOOS=linux go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv srv_arm64

build_windows_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv.exe; GOARCH=amd64 GOOS=windows go build -o srv.exe ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv.exe srv_amd64.exe

build_windows_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv.exe; GOARCH=arm64 GOOS=windows go build -o srv.exe ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv.exe srv_arm64.exe

install:
	go install

clean:
	go clean; rm -rf pkg

changes:
	./changes.sh
	./last_changes.sh

test : test_code

test_code:
	benchmark_insert.sh
	benchmark_search.sh
#     go test -test.v .

# here is an example for execution of individual test
# go test -v -run TestFilesDB
