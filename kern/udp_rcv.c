// go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf_endian.h>

#define AF_INET 2    // ipv4
#define AF_INET6 10  // ipv6

struct Udp_event {
  u16 flag;  // 0:send,1:recv
  u32 pid;
  unsigned __int128 daddr;  // 目的ip地址
  u16 dport;
  unsigned __int128 saddr;  // 源IP地址
  u16 sport;
  u16 len;
  __u8 comm[TASK_COMM_LEN];
};
const struct Udp_event *unused __attribute__((unused));
struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 1 << 24);
} events_udprcv_rb SEC(".maps");

SEC("kprobe/udp_sendmsg")
int kprobe_udp_sendmsg(struct pt_regs *ctx) {
  struct Udp_event event = {};
  event.flag = 0;
  event.pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event.comm, sizeof(event.comm));

  struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
  event.len = PT_REGS_PARM3(ctx);

  u16 family = BPF_CORE_READ(sk, __sk_common.skc_family);
  if (family == AF_INET) {
    event.daddr = BPF_CORE_READ(sk, __sk_common.skc_daddr);
    event.saddr = BPF_CORE_READ(sk, __sk_common.skc_rcv_saddr);
  } else if (family == AF_INET6) {
    struct in6_addr dtemp = BPF_CORE_READ(sk, __sk_common.skc_v6_daddr);
    struct in6_addr stemp = BPF_CORE_READ(sk, __sk_common.skc_v6_rcv_saddr);
    event.daddr = *(unsigned __int128 *)&dtemp;
    event.saddr = *(unsigned __int128 *)&stemp;
  } else {
    return 0;
  }
  event.dport = bpf_ntohs(BPF_CORE_READ(sk, __sk_common.skc_dport));
  event.sport = BPF_CORE_READ(sk, __sk_common.skc_num);
  struct Udp_event *task_info =
      bpf_ringbuf_reserve(&events_udprcv_rb, sizeof(struct Udp_event), 0);
  if (!task_info) {
    return 0;
  } else {
    *task_info = event;
    bpf_ringbuf_submit(task_info, 0);
  }
  return 0;
}

// 接受数据包长度和网络五元组二选一
SEC("kprobe/udp_recvmsg")
int kprobe_udp_recvmsg(struct pt_regs *ctx) {
  struct Udp_event event = {};
  event.flag = 1;
  event.pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event.comm, sizeof(event.comm));

  struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
  event.len = PT_REGS_PARM3(ctx);

  u16 family = BPF_CORE_READ(sk, __sk_common.skc_family);
  if (family == AF_INET) {
    event.saddr = BPF_CORE_READ(sk, __sk_common.skc_daddr);
    event.daddr = BPF_CORE_READ(sk, __sk_common.skc_rcv_saddr);
  } else if (family == AF_INET6) {
    struct in6_addr dtemp = BPF_CORE_READ(sk, __sk_common.skc_v6_daddr);
    struct in6_addr stemp = BPF_CORE_READ(sk, __sk_common.skc_v6_rcv_saddr);
    event.saddr = *(unsigned __int128 *)&dtemp;
    event.daddr = *(unsigned __int128 *)&stemp;
  } else {
    return 0;
  }
  event.sport = bpf_ntohs(BPF_CORE_READ(sk, __sk_common.skc_dport));
  event.dport = BPF_CORE_READ(sk, __sk_common.skc_num);
  struct Udp_event *task_info =
      bpf_ringbuf_reserve(&events_udprcv_rb, sizeof(struct Udp_event), 0);
  if (!task_info) {
    return 0;
  } else {
    *task_info = event;
    bpf_ringbuf_submit(task_info, 0);
  }
  return 0;
}

char LICENSE[] SEC("license") = "GPL";
