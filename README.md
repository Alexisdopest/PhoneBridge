# PhoneBridge

PhoneBridge 是一款轻量级的 Windows 与 iOS 跨设备局域网协作工具。核心理念是通过 iOS 原生“快捷指令 (Shortcuts)”实现将手机端的文本一键发送并自动写入 Windows 系统的剪贴板，以及文件/图片的静默传输保存。

## 🌟 核心特性 (v1.0 Release)
- **极速纯内网传输**：局域网直接 P2P 传输，不经外网，保护数据隐私。
- **系统托盘驻留**：无控制台黑框打扰，安静地运行在 Windows 任务栏系统托盘。
- **开机自启支持**：一键设置随系统启动，真正做到“装完即忘”。
- **设备配对鉴权**：利用 Bearer Token 进行强鉴权，防止局域网内其他设备任意调用。
- **零配置寻址 (mDNS)**：自动在局域网广播 `_phonebridge._tcp` 服务。
- **无缝剪贴板同步**：手机端执行指令后，文字立刻到达电脑剪贴板。
- **文件与图片直传**：自动分类保存在 Windows `Downloads\PhoneBridge` 目录下。

## 🚀 快速开始

### 1. Windows 服务端运行
最简单的方式是直接下载编译好的程序：

1. 在 GitHub 的 [Releases](#) 页面下载发布包或最新构建的 `phonebridge.exe`。
2. 双击运行程序。
3. **【重要】防火墙提示**：首次运行时，Windows Defender 防火墙可能会弹窗提示网络访问拦截。请务必勾选“允许访问”（专用网络和公用网络均建议勾选），否则手机端快捷指令将因为网络不通而报错。
4. 您将在右下角任务栏看到 PhoneBridge 托盘图标。右键菜单可开启“开机启动”或“打开接收文件夹”。

*(如果您希望自行编译，请克隆仓库并使用 `go build -ldflags="-H windowsgui" -o phonebridge.exe ./cmd/phonebridge` 构建。)*

### 2. 扫码极速配置 (iOS 官方快捷指令)

为了安全，PhoneBridge 会在首次启动时生成 32 字节的高强度随机配对 Token。我们为您准备了官方的快捷指令生态，**只需扫码一次，终身免配置**：

1. **核心配置 (必装)**: [PhoneBridge 配对](https://www.icloud.com/shortcuts/ea70424836cf4ed5b65287318901c693)
   *用法：在电脑右下角托盘点击“显示配对二维码”，手机运行此指令扫码，配置将安全保存在您的手机本地。*

2. **文字直传 (必装)**: [PhoneBridge 发送文本](https://www.icloud.com/shortcuts/7fc82e0fc0614fcab803dab4536a2c13)
   *用法：复制任何文字，运行此指令，文字将瞬间到达 Windows 剪贴板。*

3. **文件直传 (进阶)**: [PhoneBridge 传文件](https://www.icloud.com/shortcuts/68c2335edf054fddaa28dcd0dcfadf43)
   *用法：在相册或文件中选中资源，点击系统自带的“分享”按钮，在菜单中选择“PhoneBridge 传文件”，即可瞬间发送到电脑 `Downloads/PhoneBridge` 目录。*

---

### 3. “零点击”终极体验方案 🚀

快捷指令不应该每次都去桌面点击！配置完成后，强烈建议进行以下增强设置：

#### 玩法 A：背面轻点 (肌肉记忆)
1. 在 iPhone 中进入 `设置 -> 辅助功能 -> 触控 -> 轻点背面`。
2. 将 `轻点两下` 绑定为您刚刚创建的 PhoneBridge 发送快捷指令。
3. **体验**：在手机上选中文字点击“复制”，然后**食指在手机背面敲两下**，电脑上就可以直接 `Ctrl+V` 了！

#### 玩法 B：分享菜单直传 (免复制)
1. 在快捷指令设置中，开启 **“在共享表单中显示”**。
2. **体验**：连“复制”这一步都省了！在任何 App 里划选文字、或者在相册看图时，直接点系统的“分享”按钮，点击 PhoneBridge 图标瞬间直传电脑！

#### 玩法 C：操作按钮 (iPhone 15 Pro 及以上)
将手机侧边的“操作按钮”绑定为 PhoneBridge，实现物理按键级的外挂体验。

## 📄 架构说明
项目遵循高内聚设计，包含以下核心模块：
- `internal/app`: 全局生命周期与托盘守护。
- `internal/server`: HTTP 服务及 WebSocket 预留。
- `internal/auth`: Bearer Token 请求拦截校验。
- `internal/clipboard`: 系统剪贴板底层互斥封装。
- `internal/storage`: 防碰撞文件落地存储。
- `internal/discovery`: 局域网 mDNS 广播。
- `internal/tray & autostart`: Windows 原生托盘菜单与系统注册表管理。
