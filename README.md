# PhoneBridge

PhoneBridge 是一款轻量级的 Windows 与 iOS 跨设备局域网协作工具。核心理念是通过 iOS 原生“快捷指令 (Shortcuts)”实现将手机端的文本一键发送并自动写入 Windows 系统的剪贴板，未来还将支持文件的静默传输。

## 🌟 核心特性 (v1.0 Milestone 1)
- **极速纯内网传输**：局域网直接 P2P 传输，不经外网，保护数据隐私。
- **设备配对鉴权**：利用 Bearer Token 进行强鉴权，防止局域网内其他设备任意调用。
- **无缝剪贴板同步**：手机端执行指令后，文字立即出现在电脑剪贴板中。

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
启动成功后，控制台会输出当前监听的端口和静态配对凭证（默认为 `123456`，后续版本将加入随机生成与界面）。

### 2. iOS 快捷指令配置

为了让您的 iPhone 能够与电脑通信，请在“快捷指令” App 中创建一个新指令，按照以下步骤设置：

1. **添加操作**：“文本”（可以随意输入一些测试内容，例如 `Hello World`），或使用“获取剪贴板”、“快捷指令输入”。
2. **添加操作**：“获取 URL 内容” (Get contents of URL)
   - **URL**: `http://<您的电脑局域网IP>:8080/api/clipboard` (请替换为您电脑的实际局域网 IP)
   - **方法 (Method)**: `POST`
   - **头部 (Headers)**: 点击“添加新标题”
     - 键：`Authorization`
     - 值：`Bearer 123456`
   - **请求体 (Request Body)**: 选择 `文件 (File)`，并将“文件”参数选择为上方第1步产生的“文本”变量。

配置完成后，点击运行快捷指令。在您的 Windows PC 上按下 `Ctrl + V`，即可看到刚刚发送的文本内容！

## 📄 目录结构说明
```text
internal/
├── server      # HTTP 服务端引擎及 API Router
├── auth        # Bearer Token 请求拦截校验中间件
├── clipboard   # 剪贴板底层 API 及互斥锁封装
└── utils       # 工具函数 (Token生成等)
cmd/
└── phonebridge # 守护程序入口主文件
```
