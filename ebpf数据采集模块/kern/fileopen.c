// Copyright (C) 2024 Tianzhenxiong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

#include "common.h"

struct fileopen_event {
  u64 timestamp;
  u32 pid;
  u8 comm[TASK_COMM_LEN];
  u8 filename[FILE_MAXDEPTH][FILEPATH_LEN];
  u8 fsname[FSNAME_LEN];
  u8 cid[CONTAINER_ID_LEN];
};
const struct fileopen_event *unused_ __attribute__((unused));

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 1 << 24);
} fileopen_rb SEC(".maps");

static int add_head_slash(char *str) {
  if (str[0] == '/' && str[1] == 0) {
    char empty_str[FILEPATH_LEN] = "";
    bpf_probe_read_kernel_str(str, FILEPATH_LEN, empty_str);
    return -1;
  }
  char tmp[FILEPATH_LEN];
  bpf_probe_read_kernel_str(tmp, FILEPATH_LEN, str);
  char *_str = &str[1];
  bpf_probe_read_kernel_str(_str, FILEPATH_LEN - 1, tmp);
  str[0] = '/';
  return 1;
}

static void get_dentry_name_core(struct dentry *den, char *name) {
  u8 *namep = (u8 *)BPF_CORE_READ(den, d_name.name);
  bpf_core_read_str(name, FILEPATH_LEN, namep);
  add_head_slash(name);
}

// https://elixir.bootlin.com/linux/v6.5/source/fs/namei.c#L3812
SEC("kretprobe/do_filp_open")
int BPF_KRETPROBE(kretprobe_do_filp_open, struct file *filp) {
  struct task_struct *cur_task = (struct task_struct *)bpf_get_current_task();
  if (get_task_level_core(cur_task) == 0) {
    return 0;
  }
  struct fileopen_event *event =
      bpf_ringbuf_reserve(&fileopen_rb, sizeof(struct fileopen_event), 0);
  if (!event) {
    return 0;
  }
  event->timestamp = bpf_ktime_get_ns();
  event->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(event->comm, sizeof(event->comm));
  get_cid_core(cur_task, event->cid);
  struct file *fi = filp;
  u8 *fsnamep = (u8 *)BPF_CORE_READ(fi, f_inode, i_sb, s_type, name);
  bpf_core_read_str(&event->fsname, FSNAME_LEN, fsnamep);

  struct dentry *cur_dentry = (struct dentry *)BPF_CORE_READ(fi, f_path.dentry);
  int offset = 0;
  for (int i = 0; i < FILE_MAXDEPTH; i++) {
    if (cur_dentry == NULL) break;
    get_dentry_name_core(cur_dentry, event->filename[i]);
    if (event->filename[i][0] == 0) break;

    cur_dentry = (struct dentry *)BPF_CORE_READ(cur_dentry, d_parent);
  }
  if (event->filename[0][0] == 0) {
    bpf_ringbuf_discard(event, 0);
  } else {
    bpf_ringbuf_submit(event, 0);
  }
  return 0;
}