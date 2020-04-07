package usecases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/models"
	"server/repositories"
	"server/utils"
	"strconv"
	"strings"
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
			Data:    err.Error(),
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
		//设定上传文件的大小 单位MB
		f.Size, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", frec.Size/1024/1024), 64)
		//设定文件上传时间
		f.UploadTime = utils.Unix2DateTime(frec.Uptime)
		//设定哈希值 [重要]
		f.HashCode = frec.HashCode
		//设定文件在集群中的存储路径
		f.FilePath = frec.FilePath
		//设置 MIME
		if frec.Mime != "" {
			f.Mime = frec.Mime
		}
		//设置fsk fsk = username_base62(unix_stamp)
		f.FSK = username + "_" + utils.Unix2Base62()
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
	fid, _ := primitive.ObjectIDFromHex(id)
	var file = &models.File{
		Pid:        fid,
		Id:         primitive.NewObjectID(),
		FileName:   dirname,
		UploadTime: utils.GetNowDateTime(),
		IsDir:      true,
		IsShare:    false,
		Privacy:    false,
		FSK:        username + "_" + utils.Unix2Base62(),
	}
	if err := fr.CreateDir(username, dirname, file); err != nil {
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
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid, _ := primitive.ObjectIDFromHex(id)
	if err := fr.DeleteDir(username, fid); err != nil {
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
	fid, _ := primitive.ObjectIDFromHex(id)
	if err := fr.RenameFile(username, dirname, fid); err != nil {
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

func (this *FileUC) MoveFile(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	pid := ctx.Query("pid")
	if username == "" || id == "" || pid == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid, _ := primitive.ObjectIDFromHex(id)
	fpid, _ := primitive.ObjectIDFromHex(pid)
	if err := fr.MoveFile(username, fid, fpid); err != nil {
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

func (this *FileUC) ShareList(ctx *gin.Context) (int, *List) {
	username := ctx.Query("username")
	if username == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if list, err := fr.ShareList(username); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "获取共享列表失败",
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

func (this *FileUC) OTTHShareFile(ctx *gin.Context) (int, *List) {
	key := ctx.Query("key")
	if key == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	//解析key=> fsk=> 截取 username_base62
	//args[0] 为用户名 args[1] 为unix的base62编码
	var args = strings.Split(key, "_")
	username := args[0]
	//先获取文件信息
	f, err := fr.ShareFileInfo(username, key)
	if err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "未找到共享文件",
			Data:    err.Error(),
		}
	}
	//不是目录 直接返回该文件信息
	if !f.IsDir {
		return StatusOK, &List{
			Code:    StatusOK,
			Message: f.FileName,
			Data:    []models.File{*f},
			Count:   1,
		}
	}
	//是目录 返回该目录信息
	if list, err := fr.OTTHShareFile(username, f.Id); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "获取共享列表失败",
			Data:    err.Error(),
		}
	} else {
		return StatusOK, &List{
			Code:    StatusOK,
			Message: f.FileName,
			Data:    list,
			Count:   len(list),
		}
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
	fid, _ := primitive.ObjectIDFromHex(id)
	fsk := username + "_" + utils.Unix2Base62()
	if err := fr.ShareFile(username, fid, fsk); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "数据共享失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) CancelShare(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid, _ := primitive.ObjectIDFromHex(id)
	if err := fr.CancelShare(username, fid); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "取消共享失败",
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

func (this *FileUC) ListSecret(ctx *gin.Context) (int, *List) {
	username := ctx.Query("username")
	if username == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if list, id, err := fr.ListSecret(username); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "获取Secret文件夹数据失败",
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

func (this *FileUC) UsageRate(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	if username == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if fs, err := fr.UsageRate(username); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "查找数据失败",
			Data:    err.Error(),
		}
	} else {
		return StatusOK, &Response{
			Code:    StatusOK,
			Message: "ok",
			Data:    fs,
		}
	}
}

func (this *FileUC) CollectionList(ctx *gin.Context) (int, *List) {
	username := ctx.Query("username")
	if username == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if list, id, err := fr.CollectionList(username); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseDelete,
			Message: "获取收藏列表失败",
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

func (this *FileUC) CancelCollection(ctx *gin.Context) (int, *Response) {
	username := ctx.Query("username")
	id := ctx.Query("id")
	if username == "" || id == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	fid, _ := primitive.ObjectIDFromHex(id)
	if err := fr.CancelCollection(username, fid); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "取消收藏失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) CollectionFile(ctx *gin.Context) (int, *Response) {
	var fcr = new(models.FileCollectionRecv)
	if err := ctx.Bind(fcr); err != nil {
		return StatusClientError, &Response{
			Code:    ErrorParameterParse,
			Message: "解析参数错误",
		}
	}
	if fcr.Username == "" || fcr.FSK == "" || fcr.Filename == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	//获取收藏目录的ID
	ff, err := fr.FindFileByFilename(fcr.Username, "/@")
	if err != nil {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "系统错误！请联系管理员排查",
		}
	}
	fc := &models.FileCollection{
		Pid:            ff.Id,
		Id:             primitive.NewObjectID(),
		Filename:       fcr.Filename,
		FSK:            fcr.FSK,
		CollectionTime: utils.GetNowDateTime(),
	}
	if err := fr.CollectionFile(fcr.Username, fc); err != nil {
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "数据收藏失败",
			Data:    err.Error(),
		}
	}
	return StatusOK, &Response{
		Code:    StatusOK,
		Message: "ok",
	}
}

func (this *FileUC) UserOTTHShareList(ctx *gin.Context) (int, *List) {
	args := ctx.Query("args")
	if args == "" {
		return StatusClientError, &List{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	if u, err := ur.FindOneByUsernameOrSid(args); err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseQuery,
			Message: err,
		}
	} else {
		if list, err := fr.ShareList(u.Username); err != nil {
			return StatusServerError, &List{
				Code:    ErrorDatabaseQuery,
				Message: err,
			}
		} else {
			return StatusOK, &List{
				Code:    StatusOK,
				Count:   len(list),
				Message: "ok",
				Data:    list,
			}
		}
	}
}
