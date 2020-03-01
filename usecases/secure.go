package usecases

import (
	"github.com/gin-gonic/gin"
	"server/utils"
)

// secure 安全相关的业务处理

type SecureUC struct {
}

func NewSecureUC() SecureInterface {
	return &SecureUC{}
}

func (s SecureUC) GetToken(ctx *gin.Context) (int, *Response) {
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "",
		Data:    utils.NewSecretToken(1),
	}
}

func (s SecureUC) ParseToken(ctx *gin.Context) (int, *Response) {
	st := ctx.Query("token")
	uid, err := utils.ParseSecretToken(st)
	if err != nil {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "",
		Data:    uid,
	}
}
