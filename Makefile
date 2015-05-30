# Whether to build debug
GO_BINDATA_DEBUG=true

.PHONY: build debug release all clean test install
all: debug

debug: GO_BINDATA_DEBUG=true
debug: build test

release: GO_BINDATA_DEBUG=false
release: build test

build:
	sh scripts/do.sh build

install:
	sh scripts/do.sh install

test:
	go test -cover ./...

clean:
	rm data/* plumb && \
	rmdir data && \
	go clean
