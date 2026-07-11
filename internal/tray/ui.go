package tray

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Alexisdopest/PhoneBridge/internal/pairing"
	"github.com/skip2/go-qrcode"
)

func ShowPairingQR(port string, token string) {
	data := pairing.GeneratePayload(port, token)
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return
	}

	b64 := base64.StdEncoding.EncodeToString(png)
	html := fmt.Sprintf(`
		<html>
		<head><meta charset="utf-8"><title>PhoneBridge 配对</title></head>
		<body style="display:flex; flex-direction:column; align-items:center; justify-content:center; height:100vh; background:#f0f0f0; font-family:sans-serif;">
			<h2>PhoneBridge 配对二维码</h2>
			<img src="data:image/png;base64,%s" />
			<p>请使用 iPhone 官方快捷指令扫描此二维码完成极速配置</p>
		</body>
		</html>
	`, b64)

	tmpFile := filepath.Join(os.TempDir(), "phonebridge_pair.html")
	os.WriteFile(tmpFile, []byte(html), 0644)
	exec.Command("explorer", tmpFile).Start()
}
