package usecases

import (
	"github.com/gin-gonic/gin"
)

// UserInterface 是用户模块的接口
type UserInterface interface {
	//用户认证
	Auth(ctx *gin.Context) (int, *Response)
	//用户注册接口
	SignUp(ctx *gin.Context) (int, *Response)
	//用户登陆接口
	SignIn(ctx *gin.Context) (int, *Response)
	//找回密码 发送邮件功能
	Forget2SendEmail(ctx *gin.Context) (int, *Response)
	//找回密码 重置密码功能
	Forget2ResetPassword(ctx *gin.Context) (int, *Response)
	//检测验证ID是否失效
	CheckUUID(ctx *gin.Context) (int, *Response)
	//用户退出接口
	Logout(ctx *gin.Context) (int, *Response)

	//用户个人中心
	Profile(ctx *gin.Context) (int, *Response)
	//用户自定义删除接口
	Cancellation(ctx *gin.Context) (int, *Response)
	//修改用户信息接口
	ModifyInformation(ctx *gin.Context) (int, *Response)
	//查询单一用户接口
	FindOne(ctx *gin.Context) (int, *Response)
	//查询用户列表接口
	FindMany(ctx *gin.Context) (int, *List)
	//用户数量统计
	Census(ctx *gin.Context) (int, *Response)
	//管理员登录认证
	AuthAdminToken(ctx *gin.Context) (int, *Response)
	//用户禁用
	DisabledUser(ctx *gin.Context) (int, *Response)
}

// GroupInterface 组接口
type GroupInterface interface {
	//添加组
	AddGroup(ctx *gin.Context) (int, *Response)
	//修改组
	UpdateGroup(ctx *gin.Context) (int, *Response)
	//删除组
	DeleteGroup(ctx *gin.Context) (int, *Response)
	//罗列组
	GroupList(ctx *gin.Context) (int, *List)
}

// FileInterface 是文件模块的接口
type FileInterface interface {
	//上传文件
	UploadFile(ctx *gin.Context) (int, *Response)
	//新建文件夹
	CreateDir(ctx *gin.Context) (int, *Response)
	//删除文件
	DeleteFile(ctx *gin.Context) (int, *Response)
	//删除文件夹
	DeleteDir(ctx *gin.Context) (int, *Response)
	//修改文件名称
	RenameFile(ctx *gin.Context) (int, *Response)
	//共享列表
	ShareList(ctx *gin.Context) (int, *List)
	//第三方查看列表
	OTTHShareFile(ctx *gin.Context) (int, *List)
	//共享文件
	ShareFile(ctx *gin.Context) (int, *Response)
	//取消共享文件
	CancelShare(ctx *gin.Context) (int, *Response)
	//查找当前目录所有文件
	ListDir(ctx *gin.Context) (int, *List)
	//查看根目录
	ListRoot(ctx *gin.Context) (int, *List)
	//查看加密目录文件信息
	ListSecret(ctx *gin.Context) (int, *List)
	//查看某个文件的信息
	FileInfo(ctx *gin.Context) (int, *Response)
	//用户使用比率
	UsageRate(ctx *gin.Context) (int, *Response)
	//收藏列表
	CollectionList(ctx *gin.Context) (int, *List)
	//收藏文件
	CollectionFile(ctx *gin.Context) (int, *Response)
	//删除收藏
	CancelCollection(ctx *gin.Context) (int, *Response)
}

// SecureInterface 安全相关接口
type SecureInterface interface {
	//获取安全token
	GetToken(ctx *gin.Context) (int, *Response)
	//验证安全token
	ParseToken(ctx *gin.Context) (int, *Response)
}
