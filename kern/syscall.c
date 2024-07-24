// 获取容器syscall情况

#include "common.h"
#include "syscall_tailcall.h"
// Tracepoint bpf程序类型可选项：
// 参考https://mozillazg.com/2022/06/ebpf-libbpf-btf-powered-enabled-raw-tracepoint-common-questions.html
// 参考https://mozillazg.com/2022/05/ebpf-libbpf-raw-tracepoint-common-questions.html#hidraw-tracepoint-tracepoint
// 1. tracepoint: 最原始的tracepoint类型，基本原理是内核中预埋的静态hook点。
// 2. raw_tracepoint: 与tracepoint类似，但是可以自定义tracepoint的参数。
// 12的主要区别是，raw tracepoint 不会像 tracepoint 一样在传递上下文给 ebpf
// 程序时 预先处理好事件的参数（构造好相应的参数字段）， raw tracepoint ebpf
// 程序中访问的都是事件的原始参数。
// 因此，raw tracepoint 相比 tracepoint 性能通常会更好一点 (数据来自
// https://lwn.net/Articles/750569/ )
// 3. btf_raw_tracepoint:
// 与raw_tracepoint类似，但是可以使用BTF信息来解析tracepoint的参数。
// cat /boot/config-$(uname -r) | grep CONFIG_DEBUG_INFO_BTF # 查看是否支持BTF

struct {
  __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
  __type(key, uint32_t);
  __type(value, uint32_t);
  __uint(max_entries, 1024);
  __array(values, int());
} syscall_enter_tail_table SEC(".maps") = {
    .values =
        {
            [59] = &sys_enter_execve,
            [83] = &sys_enter_mkdir,
        },
};
struct {
  __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
  __type(key, uint32_t);
  __type(value, uint32_t);
  __uint(max_entries, 1024);
  __array(values, int());
} syscall_exit_tail_table SEC(".maps") = {
    .values =
        {
            // [59] = &sys_exit_execve,
            [83] = &sys_exit_mkdir,
        },
};

SEC("tp_btf/sys_enter")
int BPF_PROG(sys_enter, struct pt_regs *regs, long syscall_id) {
  struct task_struct *curr_task = (struct task_struct *)bpf_get_current_task();
  if (get_task_level_core(curr_task) == 0) {
    // level 0 means the task is in the root pid namespace
    return 0;
  }

  struct syscallcntkey key = {
      .syscall_id = syscall_id,
      .pid = bpf_get_current_pid_tgid() >> 32,
  };
  bpf_get_current_comm(key.comm, sizeof(key.comm));
  get_cid_core((struct task_struct *)bpf_get_current_task(), &key.cid);
  u64 *sys_cnt = bpf_map_lookup_elem(&syscall_cnt, &key);
  if (sys_cnt) {
    *sys_cnt += 1;
  } else {
    u64 zero = 0;
    bpf_map_update_elem(&syscall_cnt, &key, &zero, BPF_ANY);
  }

  bpf_tail_call(ctx, &syscall_enter_tail_table, syscall_id);

  return 0;
}

SEC("tp_btf/sys_exit")
int BPF_PROG(sys_exit, struct pt_regs *regs, long ret) {
  struct task_struct *curr_task = (struct task_struct *)bpf_get_current_task();
  if (get_task_level_core(curr_task) == 0) {
    // level 0 means the task is in the root pid namespace
    return 0;
  }
  // https://github.com/falcosecurity/libs/blob/master/driver/modern_bpf/helpers/extract/extract_from_kernel.h#L46
  uint32_t syscall_id = (uint32_t)regs->orig_ax;

  bpf_tail_call(ctx, &syscall_exit_tail_table, syscall_id);

  return 0;
}
