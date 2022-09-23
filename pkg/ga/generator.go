//nolint
package ga

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strings"
	"time"
)

func hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func generateSecret() (string, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, time.Now().UnixNano()/1000/30)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(base32encode(hmacSha1(buf.Bytes(), nil))), nil
}
