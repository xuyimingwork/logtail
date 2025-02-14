version := v1.0.0

format:
		goimports -w -l .
		gofmt -w .
		gofumpt -w .

license-check:
	# go install github.com/lsm-dev/license-header-checker/cmd/license-header-checker@latest
	license-header-checker -v -a -r apache-license.txt . go

check: license-check
		golangci-lint run

test:
		go test -v ./... -coverprofile=coverage.txt -covermode=atomic

static: webapi/index.html
	go run makestatic/makestatic.go

package:
	mkdir -p dist
	rm -f dist/*.zip
	cd dist && GOOS=linux go build ../cmd/logtail/logtail.go && zip logtail-$(version)-linux.zip logtail && rm -f logtail
	cd dist && GOOS=darwin go build ../cmd/logtail/logtail.go && zip logtail-$(version)-mac.zip logtail && rm -f logtail

build: static format check test package

linux-tools:
	cd dist && GOOS=linux go build ../cmd/logrecorder/logrecorder.go
	cd dist && GOOS=linux go build ../cmd/logrepeater/logrepeater.go
	cd dist && GOOS=linux go build ../cmd/dingmock/dingmock.go

local-tools:
	cd dist && go build ../cmd/logrecorder/logrecorder.go
	cd dist && go build ../cmd/logrepeater/logrepeater.go
	cd dist && go build ../cmd/dingmock/dingmock.go

install: format check test static
	go install cmd/logtail/logtail.go

