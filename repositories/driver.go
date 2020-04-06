package repositories

import (
	"context"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/models"
	"server/utils"
	"strings"
	"time"
	"xorm.io/core"
)

var (
	// mysql数据库-连接引擎
	engine *xorm.Engine
	// mongodb-连接客户端
	mgo *mongo.Client
	// redis-连接客户端
	rgo *redis.Client
	// ctx 上下文 目前配合mgo使用
	ctx context.Context
	// 配置文件句柄
	conf = utils.GetConfig()
)

// init 全局数据库初始化函数
func init() {
	//初始化MongoDB
	if err := InitMongoDB(); err != nil {
		panic("failed to connect mongodb:" + err.Error())
	}
	//初始化MySQL
	if err := InitMySQL(); err != nil {
		panic("failed to connect mysql:" + err.Error())
	}
	//初始化Redis
	if err := InitRedis(); err != nil {
		panic("failed to connect redis:" + err.Error())
	}
}

// InitMySQL 初始化连接MySQL
func InitMySQL() error {
	var err error
	if engine, err = xorm.NewEngine("mysql", conf.Mysql.DSN); err != nil {
		return err
	}
	//数据库连通性检测
	if err = engine.Ping(); err != nil {
		return err
	}
	//同步数据库结构
	if err = engine.Sync2(new(models.User), new(models.UserGroup), new(models.UserLog)); err != nil {
		return err
	}
	//用于设置最大打开的连接数，默认值为0表示不限制
	engine.SetMaxOpenConns(32)
	//SetMaxIdleConns用于设置闲置的连接数
	engine.SetMaxIdleConns(16)
	//设置本地时区
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	//是否开启调试
	if conf.Mysql.Debug {
		engine.Logger().SetLevel(core.LOG_DEBUG)
	}
	return nil
}

// InitMongoDB 初始化连接MongoDB [建立连接池?]
func InitMongoDB() error {
	ctx = context.TODO()
	var err error
	// 连接服务端
	if mgo, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongo.DSN)); err != nil {
		return err
	}
	// 服务连接测试
	if err = mgo.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

// 初始化连接Redis 默认使用
func InitRedis() error {
	var err error
	c := strings.Split(conf.Redis.DSN, "@")
	rgo = redis.NewClient(&redis.Options{
		Password: c[0],                         // password set
		DB:       utils.ParseStringToInt(c[1]), // use DB
		Addr:     c[2],                         // server address
	})
	// 连通性测试
	if _, err = rgo.Ping().Result(); err != nil {
		return err
	}
	return nil
}
