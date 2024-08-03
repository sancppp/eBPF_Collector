#include "common.h"
#include "vmlinux.h"
// copy
struct visit_key_t {
  u32 seq;
  u32 saddr;
  u32 daddr;
  u16 sport;
  u16 dport;
};

struct visit_value {
  u8 cid[CONTAINER_ID_USE_LEN];
  u8 comm[TASK_COMM_LEN];
  u8 sc_flag;
};

struct cnetwork_event {
  u64 timestamp;
  u32 pid;
  u8 comm[TASK_COMM_LEN];
  u8 cid[CONTAINER_ID_USE_LEN];
  u8 flag;    // 0:send 1:recv
  u32 daddr;  // 目的ip地址
  u16 dport;
  u32 saddr;  // 源IP地址
  u16 sport;
};
const struct cnetwork_event *unused_ __attribute__((unused));

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, struct socket *);
  __type(value, struct task_struct *);
  __uint(max_entries, 1024);
} socktable SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __type(key, struct visit_key_t);
  __type(value, struct visit_value);
  __uint(max_entries, 1024);
} visit_table SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 1 << 24);
} containernet_rb SEC(".maps");

// copy end

// https://elixir.bootlin.com/linux/v6.5/source/net/socket.c#L626
SEC("kretprobe/sock_alloc")
int BPF_KRETPROBE(kretprobe_sock_alloc_ret, struct socket *sock) {
  struct task_struct *curtask = (struct task_struct *)bpf_get_current_task();
  if (get_task_level_core(curtask) == 0) {
    return 0;
  }
  bpf_map_update_elem(&socktable, &sock, &curtask, BPF_ANY);
  return 0;
}

#define TCP_SKB_CB(__skb) ((struct tcp_skb_cb *)&((__skb)->cb[0]))

// https://elixir.bootlin.com/linux/v5.10.134/source/net/ipv4/tcp_output.c#L1240
SEC("kprobe/__tcp_transmit_skb")
int BPF_KPROBE(kprobe_tcp_transmit_skb) {
  if (get_task_level_core((struct task_struct *)bpf_get_current_task()) == 0) {
    return 0;
  }
  struct sock *sk = (struct sock *)PT_REGS_PARM1(ctx);
  struct task_struct *sock_task;
  struct task_struct **psock_task;
  struct socket *sock = BPF_CORE_READ(sk, sk_socket);
  psock_task = bpf_map_lookup_elem(&socktable, &sock);
  if (psock_task == NULL) {
    return 0;
  }
  sock_task = *psock_task;
  if (get_task_level_core(sock_task) == 0) {
    return 0;
  }
  struct cnetwork_event *event =
      bpf_ringbuf_reserve(&containernet_rb, sizeof(struct cnetwork_event), 0);
  if (event == NULL) {
    return 0;
  }

  u16 family = BPF_CORE_READ(sk, __sk_common.skc_family);
  if (family == AF_INET) {
    event->daddr = BPF_CORE_READ(sk, __sk_common.skc_daddr);
    event->saddr = BPF_CORE_READ(sk, __sk_common.skc_rcv_saddr);
  }

  event->dport = bpf_ntohs(BPF_CORE_READ(sk, __sk_common.skc_dport));
  event->sport = BPF_CORE_READ(sk, __sk_common.skc_num);

  event->timestamp = bpf_ktime_get_ns();
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  get_cid_core(sock_task, event->cid);
  event->flag = 0;  // send

  bpf_ringbuf_submit(event, 0);

  return 0;
}

// https://elixir.bootlin.com/linux/v5.10.134/source/net/ipv4/tcp_ipv4.c#L1668
SEC("kprobe/tcp_v4_do_rcv")
int BPF_KPROBE(kprobe_tcp_v4_do_rcv, struct sock *sk, struct sk_buff *skb) {
  if (get_task_level_core((struct task_struct *)bpf_get_current_task()) == 0) {
    return 0;
  }
  struct task_struct *sock_task;
  struct task_struct **psock_task;
  struct socket *sock = BPF_CORE_READ(sk, sk_socket);
  psock_task = bpf_map_lookup_elem(&socktable, &sock);
  if (psock_task == NULL) {
    return 0;
  }
  sock_task = *psock_task;
  if (get_task_level_core(sock_task) == 0) {
    return 0;
  }
  struct cnetwork_event *event =
      bpf_ringbuf_reserve(&containernet_rb, sizeof(struct cnetwork_event), 0);
  if (event == NULL) {
    return 0;
  }
  event->timestamp = bpf_ktime_get_ns();
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  get_cid_core(sock_task, event->cid);
  event->flag = 1;  // recv

  u16 family = BPF_CORE_READ(sk, __sk_common.skc_family);
  if (family == AF_INET) {
    event->daddr = BPF_CORE_READ(sk, __sk_common.skc_daddr);
    event->saddr = BPF_CORE_READ(sk, __sk_common.skc_rcv_saddr);
  }
  event->dport = bpf_ntohs(BPF_CORE_READ(sk, __sk_common.skc_dport));
  event->sport = BPF_CORE_READ(sk, __sk_common.skc_num);

  bpf_ringbuf_submit(event, 0);

  return 0;
}
#undef TCP_SKB_CB