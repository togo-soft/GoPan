package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File 通用型抽象文件结构体
type File struct {
	Pid        primitive.ObjectID `json:"pid"`      //父级ID ParentID
	Id         primitive.ObjectID `json:"id"`       //文件ID
	FileName   string             `json:"filename"` //文件名
	Size       int64              `json:"size"`     //文件大小 默认KB向上取整
	UploadTime string             `json:"uptime"`   //上传时间
	HashCode   string             `json:"hashcode"` //哈希值
	FilePath   string             `json:"filepath"` //文件路径
	IsDir      bool               `json:"isDir"`    //是否是文件夹
	IsShare    bool               `json:"isShare"`  //是否共享
	Privacy    bool               `json:"privacy"`  //是否是一个私有文件 单独加密存储
}

// FileRecv 接收前端file数据的结构体
type FileRecv struct {
	Pid      string `form:"pid"`      //父级ID ParentID
	Id       string `form:"id"`       //文件ID
	FileName string `form:"filename"` //文件名
	Size     int64  `form:"size"`     //文件大小 默认KB向上取整
	Uptime   int64  `form:"uptime"`   //存储端返回的上传时间戳
	HashCode string `form:"hashcode"` //哈希值
	FilePath string `form:"filepath"` //文件路径
}

// FileStorage 用户存储统计
type FileStorage struct {
}
