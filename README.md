# PhoneBridge

PhoneBridge 是一款轻量级的 Windows 与 iOS 跨设备局域网协作工具。核心理念是通过 iOS 原生“快捷指令 (Shortcuts)”实现将手机端的文本一键发送并自动写入 Windows 系统的剪贴板，以及文件/图片的静默传输保存。

## 🌟 核心特性 (v1.0 Milestone 2)
- **极速纯内网传输**：局域网直接 P2P 传输，不经外网，保护数据隐私。
- **设备配对鉴权**：利用 Bearer Token 进行强鉴权，防止局域网内其他设备任意调用。
- **零配置寻址 (mDNS)**：自动在局域网广播 `_phonebridge._tcp` 服务，摆脱固定 IP 限制（搭配支持解析 mDNS 的客户端使用）。
- **无缝剪贴板同步**：手机端执行指令后，文字立即出现在电脑剪贴板中。
- **文件与图片直传**：手机端选取图片或文件后无感发送，自动保存在 Windows 的 `Downloads\PhoneBridge` 目录下。

## 🚀 快速开始

### 1. Windows 服务端运行
请确保电脑已安装 [Go 1.22+](https://go.dev/)。

```bash
# 1. 克隆仓库
git clone https://github.com/Alexisdopest/PhoneBridge.git
cd PhoneBridge

# 2. 安装依赖
go mod tidy

# 3. 运行服务
go run ./cmd/phonebridge/main.go
```
启动成功后，控制台会输出您的配对凭证（默认为 `123456`）。

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
     - 新增字段：类型选 `文件`，键填 `file`，值选第1步获取的照片或文件。

运行快捷指令即可将文件快速投递到 Windows 的 `下载\PhoneBridge` 文件夹！

## 📄 目录结构说明
```text
internal/
├── discovery   # mDNS 局域网服务广播
├── server      # HTTP 服务端引擎及 API Router
├── auth        # Bearer Token 请求拦截校验中间件
├── clipboard   # 剪贴板底层 API 及互斥锁封装
├── storage     # 文件存储与防碰撞落地机制
└── utils       # 工具函数 (Token生成等)
cmd/
└── phonebridge # 守护程序入口主文件
```
