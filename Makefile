BINDIR        ?= bin

# Util targets
##############
.PHONY: all test 

all: build

test: go_test

build: go_build

# Go targets
#################

go_build:
	rm -f gohome
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w' -o gohome cmd/main.go

go_test:
	go test $$(go list ./... | grep -v '/vendor/')

go_fmt:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi

clean:
	rm -f gohome

