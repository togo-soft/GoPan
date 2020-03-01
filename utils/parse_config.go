package utils

// 该处定义了如何从根目录下获取并解析配置信息

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

// E 定义了读取配置文件信息的根结构
type E struct {
	Environments `yaml:"environments"`
}

// Environments 项目主要配置项[子项] 如果需要扩展 在这里添加结构来实现yaml的解析
type Environments struct {
	ProjectName string `yaml:"project_name"` //项目名称
	Debug       bool   `yaml:"debug"`        //是否开启debug模式
	Server      string `yaml:"server"`       //服务运行的host:port
	User        User   `yaml:"user"`         //用户相关配置
	Jwt         Jwt    `yaml:"jwt"`          //jwt签名配置
	File        File   `yaml:"file"`         //文件加密相关
	Email       Email  `yaml:"email"`        //邮件配置相关
	Mysql       Driver `yaml:"mysql"`        //mysql数据库配置
	Mongo       Driver `yaml:"mongodb"`      //mongodb数据库配置
	Redis       Driver `yaml:"redis"`        //redis数据库配置
}

// User 用户配置项
type User struct {
	ResetExpire int64  `yaml:"reset_expires"` //验证码有效期
	Salt        string `yaml:"salt"`          //程式通用盐
	AdminToken  string `yaml:"admin_token"`   //管理员用的唯一标识
}

// Driver 定义了数据库的连接信息 以 data-source-name 方式连接
type Driver struct {
	DSN   string `yaml:"dsn"`   //数据库连接的源名称
	Debug bool   `yaml:"debug"` //开启数据库debug模式
}

// Jwt 定义了jwt中间件的配置信息
type Jwt struct {
	SignKey     string `yaml:"sign_key"`     //jwt签名密钥
	SignMethod  string `yaml:"sign_method"`  //签名方案
	SignIssuer  string `yaml:"sign_issuer"`  //签名签发者
	SignSubject string `yaml:"sign_subject"` //签名主题
	SignExpires int64  `yaml:"sign_expires"` //签名有效期
}

// File 文件相关配置
type File struct {
	Salt string `yaml:"salt"` //盐值
	Iter int    `yaml:"iter"` //加密轮转次数
}

// Email 邮件配置相关
type Email struct {
	ServerHost string `yaml:"server_host"`   //邮件SMTP服务地址
	ServerPort int    `yaml:"server_port"`   //邮件SMTP端口地址
	UserEmail  string `yaml:"from_email"`    //发送者邮件地址
	Username   string `yaml:"from_user"`     //发送者昵称别名
	Password   string `yaml:"from_password"` //发送者密码
}

// conf 是一个全局的配置信息实例 项目运行只读取一次 是一个单例
var conf *E
var once sync.Once

// GetConfig 调用该方法会实例化conf 项目运行会读取一次配置文件 确保不会有多余的读取损耗
func GetConfig() *E {
	once.Do(func() {
		conf = new(E)
		yamlFile, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			//读取配置文件失败,停止执行
			panic("read config file error:" + err.Error())
		}
	})
	return conf
}
