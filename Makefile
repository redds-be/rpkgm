GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: compile

prep: fmt mod vet lint test

compile: clean
	@mkdir -p build/
	@go build -o build/

fmt:
	golines --max-len=120 --base-formatter=gofumpt -w $(GOFILES)

mod:
	go mod vendor
	go mod tidy

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all --fix ./...

test:
	go test ./...

clean:
	@rm -rf build/