GO ?= go

.PHONY: clean build run
run: clean build
	sudo ./ebpf_exporter

build: 
	$(GO) generate ./...
	$(GO) build

clean:
	rm -f ebpf_exporter
	find . -name "*_bpfeb.*" -delete
	find . -name "*_bpfel.*" -delete