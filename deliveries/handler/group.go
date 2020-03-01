package handler

import (
	"github.com/gin-gonic/gin"
	"server/usecases"
)

// uuc 用例层用户模块实例
var guc = usecases.NewGroupUC()

// AddGroup 添加组
func AddGroup(this *gin.Context) {
	this.JSON(guc.AddGroup(this))
}

// DeleteGroup 删除组
func DeleteGroup(this *gin.Context) {
	this.JSON(guc.DeleteGroup(this))
}

// UpdateGroup 修改组
func UpdateGroup(this *gin.Context) {
	this.JSON(guc.UpdateGroup(this))
}

// GroupList 组列表
func GroupList(this *gin.Context) {
	this.JSON(guc.GroupList(this))
}
