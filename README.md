# ebpf_exporter

目标：可扩展的ebpf数据采集框架。
该框架需要具有扩展能力，用户可以快速方便地添加新的数据采集类型。
接口化


参考：
- [text](https://github.com/cloudflare/ebpf_exporter)
- [text](https://github.com/gojue/ecapture)


sudo bpftrace -e 'kfunc:udp_sendmsg { printf("UDP sendmsg\n"); printf("%lu\n", args->len); }'
sudo bpftrace -e 'kfunc:udp_sendmsg { printf("UDP sendmsg\n"); printf("%lu\n", args->sk->__sk_common.skc_dport); }'