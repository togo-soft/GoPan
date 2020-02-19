package utils

import (
	"encoding/base64"
)

// Byte2Base64 返回base64编码
func Byte2Base64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// String2Base64 返回指定str的base64编码
func String2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
