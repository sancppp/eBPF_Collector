GO ?= go

.PHONY: clean build run
run: clean build
	sudo ./ebpf_exporter

build: 
	$(GO) generate ./...
	$(GO) build

clean:
	rm -f ebpf_exporter