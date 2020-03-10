package utils

import (
	"encoding/base64"
	"math"
	"strings"
	"time"
)

// Byte2Base64 返回base64编码
func Byte2Base64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// String2Base64 返回指定str的base64编码
func String2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64toString
func Base64toString(str string) string {
	result, _ := base64.StdEncoding.DecodeString(str)
	return string(result)
}

// Unix2Base62 将当前时间戳转为base62
func Unix2Base62() string {
	const Code62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const Length = 62
	ts := time.Now().Unix()
	result := make([]byte, 0)
	for ts > 0 {
		round := ts / Length
		remain := ts % Length
		result = append(result, Code62[remain])
		ts = round
	}
	return string(result)
}

// Base62toUnix base62的编码转为时间戳
func Base62toUnix(str string) int64 {
	var Dict = map[string]int64{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}
	const Length = 62
	str = strings.TrimSpace(str)
	var result int64 = 0
	for index, char := range []byte(str) {
		result += Dict[string(char)] * int64(math.Pow(Length, float64(index)))
	}
	return result
}
