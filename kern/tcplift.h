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

struct tcp_event {
  u16 flag;  // 0:send,1:recv
  u32 pid;
  unsigned __int128 daddr;  // 目的ip地址
  u16 dport;
  unsigned __int128 saddr;  // 源IP地址
  u16 sport;
  u16 len;
  __u8 comm[TASK_COMM_LEN];
};

const struct tcp_event *unused_ __attribute__((unused));

#endif /* __TCPLIFE_H */