GO ?= go

.PHONY: clean build all run build-static

all: clean build
	sudo ./ebpf_exporter

run:
	sudo ./ebpf_exporter

build: clean
	$(GO) generate ./...
	$(GO) build

build-static: clean
	$(GO) generate ./...
	CGO_ENABLED=0 $(GO) build -a -ldflags '-extldflags "-static"'

clean:
	rm -f ebpf_exporter
	find . -name "*_bpfeb.*" -delete
	find . -name "*_bpfel.*" -delete
