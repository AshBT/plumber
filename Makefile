# Whether to build debug
GO_BINDATA_DEBUG=true

.PHONY: build debug release all clean test
all: debug

debug: GO_BINDATA_DEBUG=true
debug: build test

release: GO_BINDATA_DEBUG=false
release: build test

build:
	sh scripts/build.sh

test:
	go test -cover ./...

clean:
	rm data/* plumb && \
	rmdir data && \
	go clean
