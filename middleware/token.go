package middleware

import (
	"github.com/gin-gonic/gin"
	"server/utils"
)

func AdminTokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("token")
		if auth == "" {
			context.JSON(TokenError, gin.H{
				"code":    TokenMissHeader,
				"message": "Request header has no token field",
			})
			context.Abort()
			return
		}
		var conf = utils.GetConfig()
		if auth != conf.User.AdminToken {
			context.JSON(TokenError, gin.H{
				"code":    TokenTampered,
				"message": "Token matching failed",
			})
			context.Abort()
			return
		}
	}

}
