package ft2xxr

import (
	"bytes"
)

func b0tos(b []byte) string {
	if i := bytes.IndexByte(b, 0); i != -1 {
		b = b[:i]
	}
	return string(b)
}
