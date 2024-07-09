# utun

之前搞过 [dtun](https://github.com/taoso/dtun/)，基于 DTLS 实现加密隧道。在实践
中发现，我的路由器性能太弱，而且不支持 AES 硬件加速，隧道流量大的时候系统负载会
很高。鉴于此，我放弃 DTLS 加密，转而使用普通的 XOR 实现简单混淆。虽然加密强度变
弱，但考虑到应用导本身还会使用 TLS 加密，所以通过牺牲 IP 层加密强度换取性能也算
是赚钱的买卖。

详细使用说明请参考我的博客文章 <https://taoshu.in/net/utun.html>

## 安装

```
go install github.com/taoso/utun/cmd/utun
```

## 服务端

```bash
utun -listen :4430 -key XXXX
```

## 客户端

```bash
utun -connect example.net:4430 -key XXXX
```

无论是服务端还是客户端，都需要自己为 tun 设备设置 IP 地址。utun 本身不维护任何
状态。
