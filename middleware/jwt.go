package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"server/utils"
	"strings"
	"time"
)

// 参考: https://godoc.org/github.com/dgrijalva/jwt-go
// 实践: https://www.jianshu.com/p/1f9915818992

const (
	//正常
	TokenOK = 200
	//异常
	TokenError = 400
	//头字段认证信息缺失
	TokenMissHeader = 4900

	//token超时
	TokenExpired = 4901
	//token格式错误
	TokenMalformed = 4902
	//token匹配uid失败
	TokenTampered = 4903
)

// JWT jwt的配置结构 这些选项需要从配置文件中读取
type JWT struct {
	SigningKey     []byte                 //签名密钥
	SigningMethod  *jwt.SigningMethodHMAC //签名方式
	SigningIssuer  string                 //签名发行人
	SigningSubject string                 //签名主题
	SigningExpires int64                  //签名有效期
}

// CustomClaims 是JWT在生成令牌时的某些声明
type CustomClaims struct {
	UserId   int64  `json:"user_id"`  //用户ID
	Username string `json:"username"` //用户名
	Email    string `json:"email"`    //邮件
	Phone    string `json:"phone"`    //用户手机
	jwt.StandardClaims
}

// JWTAuth JWT中间件 路由会注册该中间件
func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 读取认证头字段
		auth := context.Request.Header.Get("Authorization")
		uid := context.Request.Header.Get("uid")
		// 字段为空 请求非法
		if auth == "" || uid == "" {
			context.JSON(TokenError, gin.H{
				"code": TokenMissHeader,
				"message":  "Request header has no authorization field",
			})
			context.Abort()
			return
		}
		//切割头字段Authorization信息 kv[0]为类型 kv[1]为token
		kv := strings.Split(auth, " ")
		//认证非 Bearer 请求非法
		if len(kv) != 2 || kv[0] != "Bearer" {
			context.JSON(TokenError, gin.H{
				"code": TokenMissHeader,
				"message":  "Auth type invalid",
			})
			context.Abort()
			return
		}
		// token验证
		if kv[1] == "" {
			context.JSON(TokenError, gin.H{
				"code": TokenMissHeader,
				"message":  "token invalid",
			})
			context.Abort()
			return
		}
		//	解析token
		j := NewJWT()
		if status, err := j.ParseToken(kv[1], utils.ParseStringToInt64(uid)); err != nil {
			context.JSON(TokenError, gin.H{
				"code": status,
				"message":  err.Error(),
			})
			context.Abort()
			return
		}
	}
}

// NewJWT 生成JWT对象
func NewJWT() *JWT {
	var conf = utils.GetConfig()
	return &JWT{
		SigningKey:     []byte(conf.Jwt.SignKey),
		SigningMethod:  GetSigningMethod(conf.Jwt.SignMethod),
		SigningIssuer:  conf.Jwt.SignIssuer,
		SigningExpires: conf.Jwt.SignExpires,
		SigningSubject: conf.Jwt.SignSubject,
	}
}

// GetSigningMethod 通过读取配置文件 返回JWT签名方式
func GetSigningMethod(method string) *jwt.SigningMethodHMAC {
	switch method {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

// ParseToken 解析令牌 检测是否有问题
func (this *JWT) ParseToken(key string, uid int64) (int, error) {
	token, err := jwt.ParseWithClaims(key, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(this.SigningKey), nil
		})
	//数据校验
	if token.Valid {
		if c, ok := token.Claims.(*CustomClaims); ok {
			if c.UserId == uid {
				return TokenOK, nil
			}
		}
		return TokenTampered, errors.New("非法篡改UID")
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return TokenMalformed, errors.New("token格式错误")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return TokenExpired, errors.New("token已超时")
		}
	} else {
		return TokenMalformed, errors.New("无法处理该token:" + err.Error())
	}
	return TokenOK, nil
}

// NewToken 生成令牌
func (this *JWT) NewToken(uid int64, username, email, phone string) string {
	claims := CustomClaims{
		UserId:   uid,
		Username: username,
		Email:    email,
		Phone:    phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + this.SigningExpires,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    this.SigningIssuer,
			Subject:   this.SigningSubject,
		},
	}
	ss, err := jwt.NewWithClaims(this.SigningMethod, claims).SignedString(this.SigningKey)
	if err != nil {
		panic(err)
	}
	return ss
}

// RefreshToken 更新令牌
func (this *JWT) RefreshToken() {

}

// NeutralizationToken 令牌立即失效
func (this *JWT) NeutralizationToken() {

}
