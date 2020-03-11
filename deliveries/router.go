package deliveries

import (
	"github.com/gin-gonic/gin"
	"server/deliveries/handler"
	"server/middleware"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()
	//开启允许跨域
	router.Use(middleware.CORS())

	//通用路由
	var general = router.Group("/api/general")
	{
		//用户注册功能
		general.POST("/signup", handler.Signup)
		//用户登录功能
		general.POST("/signin", handler.Signin)
		//用户忘记密码 找回 发送邮件接口
		general.GET("/forget", handler.Forget)
		//用户忘记密码 找回 重置密码接口
		general.POST("/reset", handler.Reset)
		//验证码检测过期
		general.GET("/reset/check", handler.CheckUUID)
		//某共享列表 otth[on the third hand]:第三方的
		general.GET("/share", handler.OTTHShareList)
		//管理员登录认证
		general.POST("/admin", handler.AuthAdminToken)
	}

	//加载用户相关路由
	user(router)
	admin(router)
	//加载文件操作相关路由
	file(router)

	return router
}

// user 用户相关路由
func user(router *gin.Engine) {
	//用户路由 这部分路由需要JWT认证 使用JWT中间件
	var user = router.Group("/api/user", middleware.JWTAuth())
	{
		//用户token认证
		user.GET("/auth", handler.Auth)
		//查找用户的个人信息
		user.GET("/profile", handler.Profile)
		//用户帐号注销功能
		user.GET("/cancellation", handler.Cancellation)
		//用户退出登录功能
		user.GET("/logout", handler.Cancellation)
		//用户信息更新功能
		user.POST("/update", handler.ModifyInformation)
	}
}

// admin 管理员后台相关路由
func admin(router *gin.Engine) {
	//管理员路由
	var admin = router.Group("/api/admin", middleware.AdminTokenAuth())
	{
		//用户列表
		admin.GET("/user/list", handler.FindMany)
		//用户数量统计
		admin.GET("/user/count", handler.Census)
		//用户禁用
		admin.GET("/user/disabled", handler.DisabledUser)
		//添加组
		admin.POST("/group/add", handler.AddGroup)
		//删除组
		admin.GET("/group/delete", handler.DeleteGroup)
		//修改组
		admin.POST("/group/update", handler.UpdateGroup)
		//组列表
		admin.GET("/group/list", handler.GroupList)
	}
}

// file 文件相关路由
func file(router *gin.Engine) {
	var file = router.Group("/api/file", middleware.JWTAuth())
	{
		//上传文件
		file.POST("/upload", handler.UploadFile)
		//私密文件列表
		file.GET("/secret/list", handler.SecretList)
		// 新建文件夹
		file.GET("/createdir", handler.CreateDir)
		//删除文件
		file.GET("/delete", handler.DeleteFile)
		// 删除文件夹
		file.GET("/delete/dir", handler.DeleteDir)
		//重命名文件
		file.GET("/rename", handler.RenameFile)
		//共享文件
		file.GET("/share", handler.ShareFile)
		//取消共享
		file.GET("/share/cancel", handler.CancelShare)
		//共享列表
		file.GET("/share/list", handler.ShareList)
		//查看当前文件夹文件列表
		file.GET("/list", handler.ListDir)
		//查看根目录文件列表
		file.GET("/root", handler.ListRoot)
		//查看文件信息
		file.GET("/info", handler.FileInfo)
		//查看文件系统使用情况
		file.GET("/ratio", handler.FileSystemUsageRate)
		//收藏文件列表
		file.GET("/collection/list",handler.CollectionList)
		//收藏文件
		file.POST("/collection",handler.CollectionFile)
		//取消收藏
		file.GET("/collection/cancel",handler.CancelCollection)
	}
}

// secure 安全相关路由 暂时先不设置 jwt认证
func secure(router *gin.Engine) {
	var secure = router.Group("/api/secure")
	{
		secure.GET("/get", handler.GetSecretToken)
		secure.GET("/parse", handler.ParseSecretToken)
	}
}
