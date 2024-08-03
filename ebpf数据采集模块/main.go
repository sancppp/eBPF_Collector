package main

import (
	"ebpf_exporter/cmd"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go http.ListenAndServe("0.0.0.0:6060", nil)
	cmd.Execute()
}
