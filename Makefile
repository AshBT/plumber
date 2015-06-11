.PHONY: build all clean test install
all: build test

bindata/bindata.go: manager/* templates/*
	mkdir -p bindata
	go-bindata -pkg="bindata" -o=$@ manager templates

build: bindata/bindata.go
	./scripts/do.sh build

install: bindata/bindata.go
	./scripts/do.sh install

test: build
	alias gcloud=true
	alias kubectl=true
	go test -v -bench=. -cover ./...

clean:
	@rm bindata/* plumber
	@rmdir bindata
	go clean
