// go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf_endian.h>
#include "tcplift.h"

#define MAX_ENTRIES 10240
#define AF_INET 2
#define AF_INET6 10

const volatile bool filter_sport = false;
const volatile bool filter_dport = false;
const volatile __u16 target_sports[MAX_PORTS] = {};
const volatile __u16 target_dports[MAX_PORTS] = {};
const volatile pid_t target_pid = 0;
const volatile __u16 target_family = 0;

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 1 << 24);
} events_tcp_rb SEC(".maps");

SEC("kprobe/tcp_sendmsg")
int kprobe__tcp_sendmsg(struct pt_regs *ctx) {
  struct tcp_event event = {};
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
  struct tcp_event *task_info =
      bpf_ringbuf_reserve(&events_tcp_rb, sizeof(struct tcp_event), 0);
  if (!task_info) {
    return 0;
  } else {
    *task_info = event;
    bpf_ringbuf_submit(task_info, 0);
  }
  return 0;
}

SEC("kprobe/tcp_cleanup_rbuf")
int kprobe__tcp_cleanup_rbuf(struct pt_regs *ctx) {
  struct tcp_event event = {};
  event.flag = 1;
  event.pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event.comm, sizeof(event.comm));
  struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
  event.len = PT_REGS_PARM2(ctx);

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
  struct tcp_event *task_info =
      bpf_ringbuf_reserve(&events_tcp_rb, sizeof(struct tcp_event), 0);
  if (!task_info) {
    return 0;
  } else {
    *task_info = event;
    bpf_ringbuf_submit(task_info, 0);
  }
  return 0;
}

char LICENSE[] SEC("license") = "GPL";