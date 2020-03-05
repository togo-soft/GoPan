package usecases

import (
	"github.com/gin-gonic/gin"
	"server/models"
	"server/repositories"
	"server/utils"
)

// UserUC 是用户实例层的一个结构 用来实现用户接口 UserInterface
type GroupUC struct {
}

// ur 是仓库存储层的一个实例
var gr = repositories.NewGroupRepo()

// NewUserUC 会返回实例层用户模块的实例
func NewGroupUC() GroupInterface {
	return &GroupUC{}
}

func (this *GroupUC) AddGroup(ctx *gin.Context) (int, *Response) {
	var group = new(models.UserGroup)
	//检测注册信息是否解析完成
	if err := ctx.Bind(group); err != nil {
		return StatusClientError, &Response{
			Code:    ErrorParameterParse,
			Message: "解析参数错误",
		}
	}
	//判断数据是否完整
	if group.Name == "" || group.Rule == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterDefect,
			Message: "参数缺失",
		}
	}
	//写入数据库
	if id, err := gr.AddGroup(group); err != nil {
		//插入数据库失败
		return StatusServerError, &Response{
			Code:    ErrorDatabaseInsert,
			Message: "数据库插入出错",
			Data:    err,
		}
	} else {
		//操作成功
		return StatusOK, &Response{
			Code:    StatusOK,
			Message: "ok",
			Data:    utils.ParseInt64ToString(id),
		}
	}
}

func (this *GroupUC) UpdateGroup(ctx *gin.Context) (int, *Response) {
	var group = &models.UserGroup{}
	err := ctx.Bind(group)
	//检测参数是否正常
	if err != nil || group.Id == 0 {
		return StatusClientError, &Response{
			Code:    ErrorParameterParse,
			Message: "解析参数错误",
			Data:    err,
		}
	}
	if id, err := gr.UpdateGroup(group); err != nil {
		//插入数据库失败
		return StatusServerError, &Response{
			Code:    ErrorDatabaseUpdate,
			Message: "数据库操作出错",
		}
	} else {
		//操作成功
		return StatusOK, &Response{
			Code:    StatusOK,
			Message: utils.ParseInt64ToString(id),
		}
	}
}

func (this *GroupUC) DeleteGroup(ctx *gin.Context) (int, *Response) {
	qid := ctx.Query("id")
	if qid == "" {
		return StatusClientError, &Response{
			Code:    ErrorParameterParse,
			Message: "解析参数错误",
		}
	}
	id := utils.ParseStringToInt64(qid)
	if id, err := gr.DeleteGroup(id); err != nil {
		//删除数据库操作失败
		return StatusServerError, &Response{
			Code:    ErrorDatabaseDelete,
			Message: "数据库删除操作出错",
		}
	} else {
		//操作成功
		return StatusOK, &Response{
			Code:    StatusOK,
			Message: utils.ParseInt64ToString(id),
		}
	}
}

func (this *GroupUC) GroupList(ctx *gin.Context) (int, *List) {
	data, err := gr.GroupList()
	if err != nil {
		return StatusServerError, &List{
			Code:    ErrorDatabaseQuery,
			Message: "数据库查询出错",
		}
	}
	return StatusOK, &List{
		Code:    StatusOK,
		Message: "ok",
		Data:    data,
	}
}
