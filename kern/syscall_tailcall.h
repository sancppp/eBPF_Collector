#ifndef __BPF_SYSCALL_TAILCALL_H
#define __BPF_SYSCALL_TAILCALL_H

struct syscall_event {
  u8 flag;  // 0 for enter, 1 for exit, 2 counter
  u32 pid;
  u8 comm[TASK_COMM_LEN];
  u32 syscall_id;
  u64 timestamp;
  long ret;
  u8 cid[CONTAINER_ID_LEN];
  u8 info[INFO_LEN];  // extra info
};
const struct syscall_event *unused_ __attribute__((unused));

// event ringbuf
struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 1 << 24);
} syscall_rb SEC(".maps");

// syscall_idcntmap
struct syscallcntkey {
  u8 cid[CONTAINER_ID_LEN];
  u32 pid;
  u8 comm[TASK_COMM_LEN];
  u32 syscall_id;
};
const struct syscallcntkey *unused__ __attribute__((unused));

struct {
  __uint(type, BPF_MAP_TYPE_LRU_HASH);
  __type(key, struct syscallcntkey);
  __type(value, u64);
  __uint(max_entries, 1024);
} syscall_cnt SEC(".maps");

SEC("tp_btf/sys_enter")
int BPF_PROG(sys_enter_mkdir, struct pt_regs *regs) {
  // test
  // uint32_t mode = PT_REGS_PARM2_CORE(regs);
  // bpf_printk("mkdir mode: %d\n", mode);

  struct syscall_event *event =
      bpf_ringbuf_reserve(&syscall_rb, sizeof(struct syscall_event), 0);
  if (!event) {
    return 0;
  }
  event->flag = 0;
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  event->syscall_id = (uint32_t)regs->orig_ax;
  event->timestamp = bpf_ktime_get_ns();
  get_cid_core((struct task_struct *)bpf_get_current_task(), event->cid);

  bpf_ringbuf_submit(event, 0);

  return 0;
}

SEC("tp_btf/sys_exit")
int BPF_PROG(sys_exit_mkdir, struct pt_regs *regs, long ret) {
  // bpf_printk("mkdir ret: %d\n", ret);
  struct syscall_event *event =
      bpf_ringbuf_reserve(&syscall_rb, sizeof(struct syscall_event), 0);
  if (!event) {
    return 0;
  }
  event->flag = 1;
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  event->syscall_id = (uint32_t)regs->orig_ax;
  event->ret = ret;
  event->timestamp = bpf_ktime_get_ns();
  get_cid_core((struct task_struct *)bpf_get_current_task(), event->cid);

  bpf_ringbuf_submit(event, 0);
  return 0;
}

#define __user
SEC("tp_btf/sys_enter")
int BPF_PROG(sys_enter_execve, struct pt_regs *regs) {
  struct syscall_event *event =
      bpf_ringbuf_reserve(&syscall_rb, sizeof(struct syscall_event), 0);
  if (!event) {
    return 0;
  }
  event->flag = 0;
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  event->syscall_id = (uint32_t)regs->orig_ax;
  event->timestamp = bpf_ktime_get_ns();
  get_cid_core((struct task_struct *)bpf_get_current_task(), event->cid);
  u8 *filename = (u8 *)PT_REGS_PARM1_CORE(regs);
  bpf_core_read_user_str(&event->info, INFO_LEN, filename);

  bpf_ringbuf_submit(event, 0);

  return 0;
}
#undef __user
#endif