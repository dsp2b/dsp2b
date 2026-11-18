package blueprint

import (
	"fmt"
	"net/url"
	"strings"
)

// escape 等价于游戏的 Uri.EscapeDataString (RFC 3986):
// 仅保留非保留字符 A-Za-z0-9-_.~, 其余按字节转为大写 %XX。
// 注意不能用 url.QueryEscape(空格会变成"+"且不转义逗号等), 也不能用
// url.PathEscape(逗号/+等路径安全字符不转义, 会破坏头部逗号切分)。
func escape(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if isUnreserved(r) {
			sb.WriteRune(r)
			continue
		}
		for _, b := range []byte(string(r)) {
			fmt.Fprintf(&sb, "%%%02X", b)
		}
	}
	return sb.String()
}

// unescape 等价于游戏的 Uri.UnescapeDataString: 解码 %XX, '+' 保持原样。
// 例外: 站点旧版 Rename 使用 QueryEscape(空格->'+') 产生过含裸 '+' 的存量蓝图,
// 而游戏的 EscapeDataString 永远不会产生裸 '+'(字面加号编码为 %2B), 因此裸 '+'
// 只会是历史改名产物, 此时按空格解码(QueryUnescape)以还原作者本意。
// 非法的转义序列原样返回(.NET 行为, Go 的 PathUnescape 会报错)。
func unescape(s string) string {
	if strings.Contains(s, "+") {
		if u, err := url.QueryUnescape(s); err == nil {
			return u
		}
		return s
	}
	if u, err := url.PathUnescape(s); err == nil {
		return u
	}
	return s
}

func isUnreserved(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z',
		r >= 'a' && r <= 'z',
		r >= '0' && r <= '9',
		r == '-', r == '_', r == '.', r == '~':
		return true
	}
	return false
}
