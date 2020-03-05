package repositories

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"server/models"
	"server/utils"
	"time"
)

// UserRepo 操作用户模型 实现了UserRepoInterface 接口
type UserRepo struct {
}

// NewUserRepo 返回一个UserRepo对象
func NewUserRepo() UserRepoInterface {
	return &UserRepo{}
}

// InitFS 注册用户后 初始化虚拟文件系统
func (this *UserRepo) InitFS(username string) error {
	collection := mgo.Database("file").Collection(username)
	// 建立根目录
	f := &models.File{
		Id:         primitive.NewObjectID(),
		FileName:   "/",
		IsDir:      true,
		UploadTime: utils.GetNowDateTime(),
	}
	if _, err := collection.InsertOne(ctx, f); err != nil {
		return err
	}
	// 建立私密空间
	s := &models.File{
		Id:         primitive.NewObjectID(),
		FileName:   "/#",
		IsDir:      true,
		UploadTime: utils.GetNowDateTime(),
	}
	_, err := collection.InsertOne(ctx, s)
	return err
}

// InitFileStorage 用户存储统计
func (this *UserRepo) InitFileStorage(username string) error {
	collection := mgo.Database("file").Collection(username)
	f := &models.FileStorage{
		Username:  username,
		UsedSize:  "",
		TotalSize: "512",
	}
	_, err := collection.InsertOne(ctx, f)
	return err
}

// Insert 将user信息插入数据库
func (this *UserRepo) Insert(u *models.User) (int64, error) {
	return engine.Insert(u)
}

// Delete 根据ID删除user表信息
func (this *UserRepo) Delete(id int64) (int64, error) {
	return engine.Delete(&models.User{Id: id})
}

// Update 更新user信息
func (this *UserRepo) Update(u *models.User) (int64, error) {
	return engine.Id(u.Id).Update(u)
}

// FindOneById 根据ID查询一个user信息
func (this *UserRepo) FindOneById(id int64) *models.User {
	var u = &models.User{}
	if _, err := engine.ID(id).Get(u); err != nil {
		log.Fatal("find user error:", err)
	}
	return u
}

// FindOneByField 根据已有字段 查询用户信息
func (this *UserRepo) FindOneByField(u *models.User) *models.User {
	_, err := engine.Where("username = ?", u.Username).Or("email = ?", u.Email).Or("phone = ?", u.Phone).Get(u)
	if err != nil {
		log.Fatal("find user error:", err)
	}
	return u
}

// RecordExist 根据用户登陆时的名称和密码 检测user信息是否存在 并返回user信息和是否存在记录
func (this *UserRepo) LoginExist(name, password string) (bool, *models.User) {
	var u = new(models.User)
	has, _ := engine.Where("username = ?", name).Or("email = ?", name).Or("phone = ?", name).And("password = ?", password).Get(u)
	return has, u
}

// RecordExist 检测拥有该属性的对象是否存在
func (this *UserRepo) RecordExist(u *models.User) (bool, error) {
	return engine.Exist(u)
}

// RecordGet 检测是否存在并返回该对象
func (this *UserRepo) RecordGet(u *models.User) (bool, *models.User) {
	has, _ := engine.Get(u)
	return has, u
}

// 忘记密码时 发送邮件 将验证参数写入redis 设置过期时间1h
func (this *UserRepo) SetResetArgs(uuid string, uid, expires int64) error {
	//过期时间为0 则强制为3600
	if expires == 0 {
		expires = 3600
	}
	err := rgo.Set(uuid, uid, time.Duration(expires)*time.Second).Err()
	if err != nil {
		panic("redis:" + err.Error())
	}
	return err
}

// GetResetArgs 重置忘记密码时 获得参数 检验是否有效
func (this *UserRepo) GetResetArgs(uuid string) (string, error) {
	return rgo.Get(uuid).Result()
}

// FindMany 用户列表
func (this *UserRepo) FindMany() ([]*models.User, error) {
	all := make([]*models.User, 0)
	err := engine.Find(&all)
	return all, err
}

// Census 统计用户表的数量
func (this *UserRepo) Census() (int64, error) {
	return engine.Count(new(models.User))
}

// DisabledUser 用户禁用
func (this *UserRepo) DisabledUser(uid int64, status bool) error {
	user := new(models.User)
	user.Status = !status
	_, err := engine.Id(uid).Cols("status").Update(user)
	log.Println("change user status:", err)
	return err
}
