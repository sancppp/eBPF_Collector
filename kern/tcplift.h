// SPDX-License-Identifier: GPL-2.0
/* Copyright (c) 2022 Hengqi Chen */
#ifndef __TCPLIFE_H
#define __TCPLIFE_H

#define MAX_PORTS 1024
#define TASK_COMM_LEN 16

struct ident {
  __u32 pid;
  char comm[TASK_COMM_LEN];
};

struct Tcplife_event {
  unsigned __int128 saddr;
  unsigned __int128 daddr;
  __u64 ts_us;
  __u64 span_us;
  __u64 rx_b;
  __u64 tx_b;
  __u32 pid;
  __u16 sport;
  __u16 dport;
  __u16 family;
  __u8 comm[TASK_COMM_LEN];
};
const struct Tcplife_event *unused __attribute__((unused));

#endif /* __TCPLIFE_H */