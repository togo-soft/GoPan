package repositories

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/models"
)

// UserRepoInterface 用户操作接口
type UserRepoInterface interface {
	//初始化文件系统
	InitFS(username string) error
	//添加用户接口
	Insert(user *models.User) (int64, error)
	//删除用户接口
	Delete(id int64) (int64, error)
	//修改用户接口
	Update(user *models.User) (int64, error)

	//根据id查询单个用户接口
	FindOneById(id int64) *models.User
	//根据字段查询用户信息接口
	FindOneByField(u *models.User) *models.User
	//登陆时记录是否存在
	LoginExist(name, password string) (bool, *models.User)
	//检测对象是否存在
	RecordExist(u *models.User) (bool, error)
	//检测对象是否存在 并返回user对象
	RecordGet(u *models.User) (bool, *models.User)
	// 忘记密码时 发送邮件 将验证参数写入redis 过期时间expires秒
	SetResetArgs(uuid string, uid, expires int64) error
	// 重置忘记密码时 获得参数 检验是否有效
	GetResetArgs(uuid string) (string, error)
	//查询用户列表接口
	FindMany() ([]*models.User, error)
	//用户数量统计
	Census() (int64, error)
	//用户禁用
	DisabledUser(int64, bool) error
}

type GroupRepoInterface interface {
	//添加组
	AddGroup(*models.UserGroup) (int64, error)
	//修改组
	UpdateGroup(*models.UserGroup) (int64, error)
	//删除组
	DeleteGroup(int64) (int64, error)
	//罗列组
	GroupList() ([]*models.UserGroup, error)
}

// FileRepoInterface 文件操作接口
type FileRepoInterface interface {
	//上传文件
	UploadFile(username string, file *models.File) error
	//新建文件夹
	CreateDir(username, dirname string, file *models.File) error
	//下载文件
	DownloadFile(username string, id primitive.ObjectID) error
	//删除文件
	DeleteFile(username string, id primitive.ObjectID) error
	//删除文件夹
	DeleteDir(username string, id primitive.ObjectID) error
	//修改文件名称
	RenameFile(username, filename string, id primitive.ObjectID) error
	//共享列表
	ShareList(username string) ([]models.File, error)
	//第三方查看共享文件信息
	OTTHShareFile(username string, pid primitive.ObjectID) ([]models.File, error)
	//共享文件
	ShareFile(username string, id primitive.ObjectID, fks string) error
	//取消共享文件
	CancelShare(username string, id primitive.ObjectID) error
	//查看文件夹数据
	ListDir(username string, pid primitive.ObjectID) ([]models.File, error)
	//查看根目录数据
	ListRoot(username string) ([]models.File, primitive.ObjectID, error)
	//查看加密目录文件信息
	ListSecret(username string) ([]models.File, primitive.ObjectID, error)
	//查看某个文件的信息
	FileInfo(username string, id primitive.ObjectID) (*models.File, error)
}
