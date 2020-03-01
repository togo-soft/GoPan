package handler

import (
	"github.com/gin-gonic/gin"
	"server/usecases"
)

var suc = usecases.NewSecureUC()

func GetSecretToken(this *gin.Context) {
	this.JSON(suc.GetToken(this))
}

func ParseSecretToken(this *gin.Context) {
	this.JSON(suc.ParseToken(this))
}
