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
	//加载用户相关路由
	user(router)
	admin(router)
	//加载文件操作相关路由
	file(router)

	return router
}

// user 用户相关路由
func user(router *gin.Engine) {
	//用户通用路由
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
	}
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
func admin(router *gin.Engine) {
	//管理员路由
	var admin = router.Group("/api/admin")
	{
		//用户列表
		admin.GET("/user", handler.FindMany)
	}
}

// file 文件相关路由
func file(router *gin.Engine) {
	var file = router.Group("/api/file", middleware.JWTAuth())
	{
		//上传文件
		file.POST("/upload", handler.UploadFile)
		// 新建文件夹
		file.GET("/createdir", handler.CreateDir)
		//下载文件
		file.GET("/download", handler.DownloadFile)
		//删除文件
		file.GET("/delete", handler.DeleteFile)
		// 删除文件夹
		file.GET("/deldir", handler.DeleteDir)
		//重命名文件
		file.GET("/rename", handler.RenameFile)
		//共享文件
		file.GET("/share", handler.ShareFile)
		//查看当前文件夹文件列表
		file.GET("/list", handler.ListDir)
		//查看根目录文件列表
		file.GET("/root", handler.ListRoot)
		//查看文件信息
		file.GET("/info", handler.FileInfo)
	}
}
