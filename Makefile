.PHONY: build test bench fuzz fmt clean

build:
	go build -o jsonstat .

test:
	go test ./...

bench:
	go test -bench=. ./...

fuzz:
	go test -fuzz=Fuzz -fuzztime=10s ./...

fmt:
	go fmt ./...

clean:
	rm -f jsonstat