package utils

// token 生成策略 该处写死 不参考配置文件

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CustomClaims 是JWT在生成令牌时的某些声明
type CustomClaims struct {
	UserId        int64 `json:"user_id"`        //用户ID
	OperationTime int64 `json:"operation_time"` //操作时间
	jwt.StandardClaims
}

// NewToken 生成一个token
func NewSecretToken(uid int64) string {
	claims := CustomClaims{
		UserId:        uid,
		OperationTime: time.Now().Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "GoPan",
			Subject:   "GoPan - Secret",
		},
	}
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString("R29QYW4=")
	if err != nil {
		panic(err)
	}
	return ss
}

// ParseToken 解析一个token 返回token的关键标识信息
func ParseSecretToken(st string) (int64, error) {
	token, err := jwt.ParseWithClaims(st, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("R29QYW4="), nil
		})
	//数据校验
	if token.Valid {
		if c, ok := token.Claims.(*CustomClaims); ok {
			return c.UserId, nil
		}
		return -1, errors.New("非法篡改")
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return -1, errors.New("token格式错误")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return -1, errors.New("token已超时")
		}
		return -1, errors.New("token错误")
	} else {
		return -1, errors.New("无法处理该token:" + err.Error())
	}
}
