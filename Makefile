BINARY_NAME := go-starter
GO_VERSION := 1.25
COVERAGE_FILE := coverage.out

.PHONY: run build test test-race lint fmt vet clean coverage

run:
	go run main.go

build:
	go build -o $(BINARY_NAME) main.go

test:
	go test -v ./...

test-race:
	go test -race -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f *.test

coverage:
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: $(COVERAGE_FILE) and coverage.html"