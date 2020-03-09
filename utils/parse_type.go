package utils

import (
	"log"
	"strconv"
)

// 类型转换相关函数在这里被定义

// ParseIntToString 会将 int 转为 string
func ParseIntToString(num int) string {
	return strconv.Itoa(num)
}

// ParseInt64ToString 会将 int64 转为 string
func ParseInt64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// ParseInt64ToInt 会将 int64 转为 int
func ParseInt64ToInt(num int64) int {
	return int(num)
}

// ParseStringToInt 会将 string 转为 int
func ParseStringToInt(s string) int {
	_t, e := strconv.Atoi(s)
	if e != nil {
		log.Fatal("parse string to int error:", e)
		return 0
	}
	return _t
}

// ParseStringToInt64 会将 string 转为 int64
func ParseStringToInt64(s string) int64 {
	_t, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		log.Fatal("parse string to int64 error:", e)
		return 0
	}
	return _t
}

// ParseStringToFloat64 将 string 转 float64
func ParseStringToFloat64(s string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal("parse string to int64 error:", err)
		return 0
	}
	return res
}