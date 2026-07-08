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
请确保电脑已安装 [Go 1.22+](https://go.dev/)。

```bash
# 1. 克隆仓库
git clone https://github.com/Alexisdopest/PhoneBridge.git
cd PhoneBridge

# 2. 安装依赖
go mod tidy

# 3. 编译发布版 (隐藏控制台，后台托盘运行)
go build -ldflags="-H windowsgui" -o phonebridge.exe ./cmd/phonebridge

# 4. 启动程序
./phonebridge.exe
```
启动成功后，您将在右下角任务栏看到 PhoneBridge 图标。通过右键菜单可以开启“开机启动”或“打开接收文件夹”。

### 2. iOS 快捷指令配置

#### 场景 A：发送文本到 Windows 剪贴板
1. **添加操作**：“文本”（输入测试内容）或“获取剪贴板”。
2. **添加操作**：“获取 URL 内容” (Get contents of URL)
   - **URL**: `http://<您的电脑局域网IP>:8080/api/clipboard`
   - **方法 (Method)**: `POST`
   - **头部 (Headers)**: `Authorization` -> `Bearer 123456`
   - **请求体 (Request Body)**: 选择 `文件 (File)`，选取第1步的文本。

#### 场景 B：发送图片或文件到 Windows
1. **添加操作**：“选择照片” 或 “获取文件”。
2. **添加操作**：“获取 URL 内容” (Get contents of URL)
   - **URL**: `http://<您的电脑局域网IP>:8080/api/upload`
   - **方法 (Method)**: `POST`
   - **头部 (Headers)**: `Authorization` -> `Bearer 123456`
   - **请求体 (Request Body)**: 选择 `表单 (Form)`
     - 新增字段：类型选 `文件`，键填 `file`，值选第1步的照片或文件。

配置完成后运行快捷指令，即可体验无缝投递！

## 📄 架构说明
项目遵循高内聚设计，包含以下核心模块：
- `internal/app`: 全局生命周期与托盘守护。
- `internal/server`: HTTP 服务及 WebSocket 预留。
- `internal/auth`: Bearer Token 请求拦截校验。
- `internal/clipboard`: 系统剪贴板底层互斥封装。
- `internal/storage`: 防碰撞文件落地存储。
- `internal/discovery`: 局域网 mDNS 广播。
- `internal/tray & autostart`: Windows 原生托盘菜单与系统注册表管理。
