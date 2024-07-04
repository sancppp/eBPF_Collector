## TCP接收和发送事件 (`event.Tcplife_event`)

### **`ebpf_tcp_byte_counter`**: 记录进程的TCP事件
对于TCP事件，记录了字节数，以及与事件相关的其他元数据。

记录以下标签：

- `timestamp`: 事件发生的时间戳。
- `protocol`: 使用的协议，对于这个计数器，值始终为`"tcp"`。**TODO 后续改为ipv4 or ipv6**
- `pid`: 发生事件的进程ID。
- `flag`: 发送消息时使用的标志。// 0:send,1:recv
- `comm`: 发生事件的进程名称。
- `daddr`: 目的地址。
- `dport`: 目的端口。
- `saddr`: 源地址。
- `sport`: 源端口。

### **`ebpf_tcp_byte_sum`**: 记录进程的TCP数据的字节数

此计数器记录了通过TCP发送/接受的总字节数。它使用以下标签来提供事件的上下文：

- `protocol`: 使用的协议，对于这个计数器，值始终为`"tcp"`。**TODO 后续改为ipv4 or ipv6**
- `flag`: 发送消息时使用的标志。// 0:send,1:recv
- `pid`: 发生事件的进程ID。
- `comm`: 发生事件的进程名称。


## UDP接收和发送事件 (`event.UdpSendmsg_event`)

###  **`ebpf_udp_byte_counter`**: 记录进程的UDP事件

**TODO 目前收包的字节数不对，还在debug**

对于UDP发送消息事件。这个计数器记录了发送的UDP消息的字节数，方向，以及与事件相关的其他元数据。

对于这个计数器，会记录以下标签：

- `timestamp`: 事件发生的时间戳。
- `protocol`: 使用的协议，对于这个计数器，值始终为`"udp"`。**TODO 后续改为ipv4 or ipv6**
- `flag`: 发送消息时使用的标志。// 0:send,1:recv
- `pid`: 发生事件的进程ID。
- `comm`: 发生事件的进程名称。
- `daddr`: 目的地址。
- `dport`: 目的端口。
- `saddr`: 源地址。
- `sport`: 源端口。

### **`ebpf_udp_byte_sum`**: 记录进程的UDP数据的字节数

此计数器记录了通过UDP发送/接受的总字节数。它使用以下标签来提供事件的上下文：

- `protocol`: 使用的协议，对于这个计数器，值始终为`"udp"`。**TODO 后续改为ipv4 or ipv6**
- `flag`: 发送消息时使用的标志。// 0:send,1:recv
- `pid`: 发生事件的进程ID。
- `comm`: 发生事件的进程名称。

