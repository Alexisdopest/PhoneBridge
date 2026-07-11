# PhoneBridge v1.0.0

这是 PhoneBridge 的首个正式版发布 (Release v1.0.0)！我们实现了 iPhone 与 Windows 之间的极速局域网无缝协同工作流。

## Features
- **零点击无感同步**：iPhone 发送，Windows 自动覆写剪贴板或归档文件。
- **系统托盘守护**：完美隐藏黑框，常驻任务栏右下角，一键支持系统开机自启。
- **动态网络寻址 (mDNS)**：自动释放 `_phonebridge._tcp` 广播服务，彻底摆脱固定 IP 配置。
- **文件极速直传**：P2P 局域网大文件直接传输，不绕行外网，保护隐私且高速流转。

## Installation
1. 在下方 `Assets` 下载 `PhoneBridge-v1.0.0.zip` 并解压。
2. 在 Windows 电脑上双击运行 `phonebridge.exe`。
3. **【必读】** 如果弹出 Windows 防火墙拦截提示，请务必勾选 **“允许访问”**。
4. 参照 [README 教程](https://github.com/Alexisdopest/PhoneBridge/blob/main/README.md) 在 iPhone 的“快捷指令 (Shortcuts)” App 中配置动作。

## Requirements
- Windows 10/11
- 搭载 iOS 15 及以上系统的 iPhone (需使用系统自带的 Shortcuts App)

## Security
- API 层全面受 `Bearer Token` 中间件鉴权保护（局域网隔离防滥用）。
- HTTP 文件存储路由内置了强力的 `Path Traversal` (路径逃逸) 免疫算法，保障宿主机核心目录安全。

## Known limitations
- 目前配对 Token 暂时硬编码为 `123456`，将在 v1.1 版本中开放动态 Token 与配置化面板界面展示。
- 当前版本暂未提供 iOS 独立 App，因此所有的自动化收发行为需借助 Shortcuts 完成。
- `WebSocket` 推送核心已经完成预埋，将在后续配套的 iOS App 发布时激活双向通知能力。
