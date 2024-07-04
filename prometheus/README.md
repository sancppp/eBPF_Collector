## TCP接收和发送事件 (`event.Tcplife_event`)

### **`ebpf_tcp_rx_event`**: 记录接收到的TCP数据的字节数
对于TCP生命周期事件，记录了接收（Rx）的字节数，以及与事件相关的其他元数据。

记录以下标签：

- `timestamp`: 事件发生的时间戳。
- `protocol`: 使用的协议，对于这个计数器，值始终为`"tcp"`。**TODO 后续改为ipv4 or ipv6**
- `pid`: 发生事件的进程ID。
- `comm`: 发生事件的进程名称。
- `daddr`: 目的地址。
- `dport`: 目的端口。
- `saddr`: 源地址。
- `sport`: 源端口。

### **`ebpf_tcp_tx_event`**: 记录发送的TCP数据的字节数

对于TCP生命周期事件，记录了发送（Tx）的字节数，以及与事件相关的其他元数据。

记录以下标签：

- `timestamp`: 事件发生的时间戳。
- `protocol`: 使用的协议，对于这个计数器，值始终为`"tcp"`。**TODO 后续改为ipv4 or ipv6**
- `pid`: 发生事件的进程ID。
- `comm`: 发生事件的进程名称。
- `daddr`: 目的地址。
- `dport`: 目的端口。
- `saddr`: 源地址。
- `sport`: 源端口。



### **`ebpf_tcp_rx_byte_sum`**: TCP接收字节总和

此计数器记录了通过TCP接收的总字节数。它使用以下标签来提供事件的上下文：

- **`protocol`**: 使用的协议，对于这个计数器，值始终为`"tcp"`。
- **`pid`**: 发生事件的进程ID。
- **`comm`**: 发生事件的进程名称。

通过这些标签，可以详细了解哪个进程从哪个源地址和端口接收了多少TCP数据。

### **`ebpf_tcp_tx_byte_sum`**: TCP发送字节总和

此计数器记录了通过TCP发送的总字节数。它使用以下标签来提供事件的上下文：

- **`protocol`**: 使用的协议，对于这个计数器，值始终为`"tcp"`。
- **`pid`**: 发生事件的进程ID。
- **`comm`**: 发生事件的进程名称。


## UDP接收和发送事件 (`event.UdpSendmsg_event`)


###  **`ebpf_udp_event`**: 记录发送的UDP消息的字节数

**TODO 目前收包的字节数不对，还在debug**

对于UDP发送消息事件，有一个计数器：`udpEventCounter`。这个计数器记录了发送的UDP消息的字节数，方向，以及与事件相关的其他元数据。

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
