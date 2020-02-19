package usecases

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"server/models"
	"server/repositories"
	"server/utils"
)

// fr 文件仓库操作句柄
var fr = repositories.NewFileRepo()

// FileUC 是文件实例层的一个结构 用来实现文件操作接口 FileInterface
type FileUC struct {
}

// NewFileUC 会返回实例层文件模块的实例
func NewFileUC() FileInterface {
	return &FileUC{}
}

func (this *FileUC) UploadFile(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	var frec = new(models.FileRecv)
	var f = new(models.File)
	//检测上传信息是否绑定成功
	if err := ctx.Bind(frec); err != nil {
		return StatusClientError, &Response{
			Code:    ErrorParameterParse,
			Message: "解析参数错误",
		}
	}
	//开始逻辑
	{
		//设定上传文件的父级目录ID
		f.Pid, _ = primitive.ObjectIDFromHex(frec.Pid)
		//设定文件的ID
		f.Id = primitive.NewObjectID()
		//设定文件名
		f.FileName = frec.FileName
		//设定上传文件的大小
		f.Size = int64(math.Ceil(float64(frec.Size) / 1024))
		//设定文件上传时间
		f.UploadTime = utils.Unix2DateTime(frec.Uptime)
		//设定哈希值 [重要]
		f.HashCode = frec.HashCode
		//设定文件在集群中的存储路径
		f.FilePath = frec.FilePath
	}
	//写入数据库
	if err := fr.UploadFile(username, f); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseInsert,
			Message: "插入数据失败",
			Data:    err.Error(),
		}
	}
	//操作成功
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) CreateDir(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	dirname := ctx.Query("dirname")
	if username == "" || id == "" || dirname == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid,_ := primitive.ObjectIDFromHex(id)
	var file = &models.File{
		Pid:        fid,
		Id:         primitive.NewObjectID(),
		FileName:   dirname,
		UploadTime: utils.GetNowDateTime(),
		IsDir:      true,
		IsShare:    false,
		Privacy:    false,
	}
	if err := fr.CreateDir(username,dirname,file); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "删除数据失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) DownloadFile(ctx *gin.Context) (int, *Response) {
	panic("implement me")
}

func (this *FileUC) DeleteFile(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid, _ := primitive.ObjectIDFromHex(id)
	if err := fr.DeleteFile(username, fid); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "删除数据失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) DeleteDir(ctx *gin.Context) (int, *Response) {
	panic("implement me")
}

func (this *FileUC) RenameFile(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	dirname := ctx.Query("dirname")
	if username == "" || id == "" || dirname == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid,_ := primitive.ObjectIDFromHex(id)
	if err := fr.RenameFile(username,dirname,fid); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "删除数据失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) ShareFile(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid,_ := primitive.ObjectIDFromHex(id)
	if err := fr.ShareFile(username,fid); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "删除数据失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) ListDir(ctx *gin.Context) (int, *List) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	_id, _ := primitive.ObjectIDFromHex(id)
	if list, err := fr.ListDir(username, _id); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "删除数据失败",
			Data:    err.Error(),
		}
	} else {
		return StatusOK, &List{
			Code:    StatusOK,
			Message: "ok",
			Data:    list,
			Count:   len(list),
		}
	}
}

func (this *FileUC) ListRoot(ctx *gin.Context) (int, *List) {
	username := ctx.Query("username")
	if username == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if list, id, err := fr.ListRoot(username); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "获取ROOT数据失败",
			Data:    err.Error(),
		}
	} else {
		return StatusOK, &List{
			Code:    StatusOK,
			Message: id,
			Data:    list,
			Count:   len(list),
		}
	}
}

func (this *FileUC) FileInfo(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	_id, _ := primitive.ObjectIDFromHex(id)
	if file, err := fr.FileInfo(username, _id); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "查找数据失败",
			Data:    err.Error(),
		}
	} else {
		return StatusOK, &Response{
			Code:    StatusOK,
			Message: "ok",
			Data:    file,
		}
	}
}
