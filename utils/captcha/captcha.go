package captcha

import (
	"bytes"
	"io"
	"time"

	"github.com/dchest/captcha"
)

var (
	globalStore = captcha.NewMemoryStore(50, time.Minute*5)
	defaultLen  = 4
)

func init() {
	captcha.SetCustomStore(globalStore)
}

// New 创建新的验证码
func New() string {
	return captcha.NewLen(defaultLen)
}

// Reload 重新加载验证码
func Reload(id string) bool {
	return captcha.Reload(id)
}

// WriteImage 图片验证码
func WriteImage(w io.Writer, id string, width, height int) error {
	return captcha.WriteImage(w, id, width, height)
}

// VerifyString 验证码验证
func VerifyString(id string, digits string) bool {
	if digits == "" {
		return false
	}
	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}
	if len(ns) == 0 {
		return false
	}

	reald := globalStore.Get(id, false)
	if reald == nil {
		return false
	}
	return bytes.Equal(ns, reald)
}
