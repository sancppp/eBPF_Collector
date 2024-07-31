# ebpf_exporter

目标：可扩展的ebpf数据采集框架。
该框架需要具有扩展能力，用户可以快速方便地添加新的数据采集类型。
接口化

## 如何扩展？

内核态C代码 + cilium/ebpf用户态代码 + event消费逻辑

## 参考
- [text](https://github.com/cloudflare/ebpf_exporter)
- [text](https://github.com/gojue/ecapture)

sudo bpftrace -e 'kfunc:udp_sendmsg { printf("UDP sendmsg\n"); printf("%lu\n", args->len); }'
sudo bpftrace -e 'kfunc:udp_sendmsg { printf("UDP sendmsg\n"); printf("%lu\n", args->sk->__sk_common.skc_dport); }'


## 网络相关

### TCP

### UDP

更正：

入口udp流量从SEC("kprobe/udp_send_skb") && SEC("kprobe/ip_send_skb")获取
出口udp流量从SEC("kprobe/udp_rcv") && SEC("kprobe/__udp_enqueue_schedule_skb")获取

确认配置：
cat /boot/config-$(uname -r) | grep CONFIG_DEBUG_INFO_BTF

## libbpf安装
没安装的话：<bpf/*.h>头文件找不到

安装：
```shell
git clone --depth 1 https://github.com/libbpf/libbpf
cd src
sudo make install

or

sudo apt install libbpf-dev 
```

## bpftrace安装
