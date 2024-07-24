#pragma once  // 与#ifndef #define #endif的作用一样，防止头文件被重复引用

#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_endian.h>

// DEFINE
#define CONTAINER_ID_LEN 32
#define CONTAINER_ID_USE_LEN 12

#define INFO_LEN 64

#define AF_INET 2
#define AF_INET6 10

#define FILEPATH_LEN 128
#define FILE_MAXDEPTH 32
#define FSNAME_LEN 64

/*--------tools-----------*/
static int get_cid_core(struct task_struct *task, u8 *cid) {
  struct css_set *css = BPF_CORE_READ(task, cgroups);
  struct cgroup_subsys_state *sbs = BPF_CORE_READ(css, subsys[0]);
  struct cgroup *cg = BPF_CORE_READ(sbs, cgroup);
  struct kernfs_node *knode = BPF_CORE_READ(cg, kn);
  struct kernfs_node *pknode = BPF_CORE_READ(knode, parent);
  u8 tmp_cid[CONTAINER_ID_LEN];
  u8 *_cid;
  if (pknode != NULL) {
    u8 *aus = (u8 *)BPF_CORE_READ(knode, name);
    // BPF_CORE_READ_STR_INTO(&tmp_cid, knode, name); ??
    bpf_core_read_str(&tmp_cid, CONTAINER_ID_LEN, aus);
    if (tmp_cid[6] == '-')
      _cid = &tmp_cid[7];
    else
      _cid = (u8 *)&tmp_cid;
    bpf_core_read_str(cid, CONTAINER_ID_USE_LEN, _cid);
  }
  return sizeof(cid);
}

static int get_task_level_core(struct task_struct *task) {
  return BPF_CORE_READ(task, nsproxy, pid_ns_for_children, level);
}

/*--------tools end-----------*/
char LICENSE[] SEC("license") = "GPL";