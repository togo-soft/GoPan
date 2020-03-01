package middleware

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
	//token匹配失败
	TokenTampered = 4903
)
