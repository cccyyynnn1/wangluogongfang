package utils

import (
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetDecodedPathParam 提取并规范化路径参数（兼容已编码/未编码、正反斜杠）
// 适用于使用 /:path 路径参数的路由；会移除前导斜杠，并将正斜杠替换为反斜杠
// 返回用于文件/注册表等后端实际操作的本地路径字符串
func GetDecodedPathParam(c *gin.Context, paramName string) (string, error) {
	raw := c.Param(paramName)
	// gin 的 *path 通常会带一个前导'/'
	if strings.HasPrefix(raw, "/") {
		raw = strings.TrimPrefix(raw, "/")
	}

	// 尝试URL解码（若已编码则解码；未编码则原样返回）
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		// 解码失败则回退原始值，尽量不中断
		decoded = raw
	}

	// 将正斜杠规范为反斜杠，方便Windows路径处理
	normalized := strings.ReplaceAll(decoded, "/", "\\")

	return normalized, nil
}

// EnsureEncodedPathParam 确保将本地路径编码为URL安全的路径段
// 用于服务端拼接URL或返回给前端用作路径参数时，避免路由分段导致404
func EnsureEncodedPathParam(localPath string) string {
	if localPath == "" {
		return ""
	}
	// 统一使用反斜杠，再做路径清理（避免 .. 等）
	p := strings.ReplaceAll(localPath, "/", "\\")
	// PathEscape 会对 ':'、'\\' 等进行百分号编码，满足参数传递需求
	return url.PathEscape(path.Clean(p))
}
