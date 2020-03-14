package models

import "time"

// User 用户表结构
type User struct {
	Id         int64  `form:"id" json:"id"`
	Username   string `form:"username" json:"username" xorm:"varchar(32) notnull unique"` //用户名
	Password   string `form:"password" json:"-" xorm:"varchar(64) notnull"`               //密码
	Fk         string `form:"-" json:"fk" xorm:"varchar(32) notnull"`                     //文件密钥
	Iv         string `form:"-" json:"iv" xorm:"varchar(32) notnull"`                     //文件密钥对应的IV向量
	Ak         string `form:"-" json:"ak" xorm:"varchar(16) notnull"`                     //文件 auth_key 与底部存储有关
	Email      string `form:"email" json:"email" xorm:"varchar(128) notnull unique"`      //邮箱
	Phone      string `form:"phone" json:"phone" xorm:"varchar(11)"`                      //电话号码
	CreateTime string `form:"-" json:"create_time" xorm:"varchar(19)"`                    //注册时间
	Status     bool   `form:"status" json:"status"`                                       //账户状态
}

// Group 用户组结构
type UserGroup struct {
	Id      int64  `form:"id" json:"id" xorm:"pk autoincr"`                    //组ID
	Name    string `form:"name" json:"name" xorm:"varchar(32) notnull"`        //组名称
	Rule    string `form:"rule" json:"rule" xorm:"varchar(16) notnull"`        //组规则
	Effect  string `form:"effect" json:"effect" xorm:"varchar(16) notnull"`    //rule对应的效果 标注rule的信息
	Explain string `form:"explain" json:"explain" xorm:"varchar(128) notnull"` //组说明
}

// UserLog 用户行为日志记录
type UserLog struct {
	Id               int64     `json:"id"`                                     //主键
	Uid              int64     `json:"uid"`                                    //用户表外键
	LastTime         time.Time `json:"last_time"`                              //上次登陆时间
	CurrentTime      time.Time `json:"current_time"`                           //本次登陆时间记录
	LastIP           string    `json:"last_ip" xorm:"varchar(46)"`             //上次登陆IP
	CurrentIP        string    `json:"current_ip" xorm:"varchar(46)"`          //本次登陆IP
	LastUserAgent    string    `json:"last_user_agent" xorm:"varchar(128)"`    //上次UserAgent
	CurrentUserAgent string    `json:"current_user_agent" xorm:"varchar(128)"` //本次UserAgent
}
