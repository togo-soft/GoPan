package handler

import (
	"github.com/gin-gonic/gin"
	"server/usecases"
)

var fuc = usecases.NewFileUC()

// UploadFile 上传文件
func UploadFile(this *gin.Context) {
	this.JSON(fuc.UploadFile(this))
}

// CreateDir 新建文件夹
func CreateDir(this *gin.Context) {
	this.JSON(fuc.CreateDir(this))
}

// DeleteFile 删除文件
func DeleteFile(this *gin.Context) {
	this.JSON(fuc.DeleteFile(this))
}

// DeleteDir 删除文件
func DeleteDir(this *gin.Context) {
	this.JSON(fuc.DeleteDir(this))
}

// MoveFile 删除文件
func MoveFile(this *gin.Context) {
	this.JSON(fuc.MoveFile(this))
}

// RenameFile 重命名文件
func RenameFile(this *gin.Context) {
	this.JSON(fuc.RenameFile(this))
}

// ShareList 共享列表
func ShareList(this *gin.Context) {
	this.JSON(fuc.ShareList(this))
}

// OTTHShareList 共享列表
func OTTHShareList(this *gin.Context) {
	this.JSON(fuc.OTTHShareFile(this))
}

// ShareFile 共享文件
func ShareFile(this *gin.Context) {
	this.JSON(fuc.ShareFile(this))
}

// CancelShare 取消共享
func CancelShare(this *gin.Context) {
	this.JSON(fuc.CancelShare(this))
}

// ListDir 查看当前文件夹文件列表
func ListDir(this *gin.Context) {
	this.JSON(fuc.ListDir(this))
}

// ListDir 查看当前文件夹文件列表
func ListRoot(this *gin.Context) {
	this.JSON(fuc.ListRoot(this))
}

// FileInfo 查看文件信息
func FileInfo(this *gin.Context) {
	this.JSON(fuc.FileInfo(this))
}

// SecretList 私密文件列表
func SecretList(this *gin.Context) {
	this.JSON(fuc.ListSecret(this))
}

// FileSystemUsageRate 文件系统使用率
func FileSystemUsageRate(this *gin.Context) {
	this.JSON(fuc.UsageRate(this))
}

// CollectionList 文件收藏列表
func CollectionList(this *gin.Context) {
	this.JSON(fuc.CollectionList(this))
}

// CollectionFile 收藏文件
func CollectionFile(this *gin.Context) {
	this.JSON(fuc.CollectionFile(this))
}

// CancelCollection 取消收藏
func CancelCollection(this *gin.Context) {
	this.JSON(fuc.CancelCollection(this))
}

// UserOTTHShareList 查询用户的共享列表
func UserOTTHShareList(this *gin.Context) {
	this.JSON(fuc.UserOTTHShareList(this))
}
