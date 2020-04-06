package handler

import (
	"github.com/gin-gonic/gin"
	"server/usecases"
)

// uuc 用例层用户模块实例
var uuc = usecases.NewUserUC()

// Auth 用户TOKEN认证
func Auth(this *gin.Context) {
	this.JSON(uuc.Auth(this))
}

// Signup 用户注册
func Signup(this *gin.Context) {
	this.JSON(uuc.SignUp(this))
}

// Signin 用户登陆
func Signin(this *gin.Context) {
	this.JSON(uuc.SignIn(this))
}

// Forget 忘记密码 发送邮件功能
func Forget(this *gin.Context) {
	this.JSON(uuc.Forget2SendEmail(this))
}

// Reset 忘记密码 重置密码功能
func Reset(this *gin.Context) {
	this.JSON(uuc.Forget2ResetPassword(this))
}

// CheckUUID 忘记密码 重置密码功能
func CheckUUID(this *gin.Context) {
	this.JSON(uuc.CheckUUID(this))
}

// Logout 用户退出
func Logout(this *gin.Context) {
	this.JSON(uuc.Logout(this))
}

// Profile 用户个人中心
func Profile(this *gin.Context) {
	this.JSON(uuc.Profile(this))
}

// FindOne 查看用户信息
func FindOne(this *gin.Context) {
	this.JSON(uuc.FindOne(this))
}

// FindMany 查看用户列表
func FindMany(this *gin.Context) {
	this.JSON(uuc.FindMany(this))
}

// Cancellation 删除用户
func Cancellation(this *gin.Context) {
	this.JSON(uuc.Cancellation(this))
}

// ModifyInformation 修改用户
func ModifyInformation(this *gin.Context) {
	this.JSON(uuc.ModifyInformation(this))
}

func Census(this *gin.Context) {
	this.JSON(uuc.Census(this))
}

// AuthAdminToken 管理员认证
func AuthAdminToken(this *gin.Context) {
	this.JSON(uuc.AuthAdminToken(this))
}

// DisabledUser 用户禁用
func DisabledUser(this *gin.Context) {
	this.JSON(uuc.DisabledUser(this))
}

// UserGroupUpdate 修改用户分组
func UserGroupUpdate(this *gin.Context) {
	this.JSON(uuc.UserGroupUpdate(this))
}

// UserLog 查询日志
func UserLog(this *gin.Context) {
	this.JSON(uuc.QueryLog(this))
}

// LogList 查询日志列表
func LogList(this *gin.Context) {
	this.JSON(uuc.QueryLogList(this))
}

// AdminResetPassword 管理员重置用户密码
func AdminResetPassword(this *gin.Context) {
	this.JSON(uuc.AdminResetPassword(this))
}
