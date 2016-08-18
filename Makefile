COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date +%FT%T%z`
VERSION_NUMBER=v0.1.0
LDFLAGS=-ldflags "-X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -X main.VersionNumber=${VERSION_NUMBER}"

default: fmt test build

build:
	go build ${LDFLAGS}

fmt:
	go fmt

test:
	go test ./...

clean:
	rm -rf document-imaging
