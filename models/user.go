package models

import "time"

// User 用户表结构
type User struct {
	Id         int64     `form:"-" json:"id"`
	Username   string    `form:"username" json:"username" xorm:"varchar(32) notnull unique"` //用户名
	Nickname   string    `form:"nickname" json:"nickname" xorm:"varchar(32) notnull"`        //昵称
	Password   string    `form:"password" json:"-" xorm:"varchar(64) notnull"`               //密码
	Fk         string    `form:"-" json:"fk" xorm:"varchar(44) notnull"`                     //文件密钥
	Ak         string    `form:"-" json:"ak" xorm:"varchar(36) notnull"`                     //文件 auth_key 与底部存储有关
	Email      string    `form:"email" json:"email" xorm:"varchar(128) notnull unique"`      //邮箱
	Phone      string    `form:"phone" json:"phone" xorm:"varchar(11)"`                      //电话号码
	CreateTime time.Time `form:"-" json:"create_time" xorm:"created"`                        //注册时间
	Status     bool      `form:"status" json:"status"`                                       //账户状态
}

// Group 用户组结构
type Group struct {
	Id   int    `form:"-" json:"id" xorm:"pk autoincr"`
	Name string `form:"group_name" json:"group_name" xorm:"varchar(32) notnull"`
	Role string `form:"role" json:"role" xorm:"varchar(32) notnull"`
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
