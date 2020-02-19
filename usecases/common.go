package usecases

// 与交付层对接的公共操作
// http状态只有为200(成功)/400(客户端错误)/500(服务端错误)
const (
	//成功
	StatusOK = 200
	//客户端错误
	StatusClientError = 400
	//服务端失败
	StatusServerError = 500

	//参数相关错误
	//参数解析错误 [来自客户端的参数解析出错]
	ErrorParameterParse = 4001
	//必要的参数缺失
	ErrorParameterDefect = 4002

	//获取错误
	//文件读取错误
	ErrorReadLocal = 5001
	//读取远程数据出错
	ErrorReadRemote = 5002
	//解析远程数据出错
	ErrorParseRemote = 5003

	//数据库相关错误
	//数据库连接出错
	ErrorDatabaseConnection = 5010
	//数据库查询错误
	ErrorDatabaseQuery = 5011
	//数据库插入出错
	ErrorDatabaseInsert = 5012
	//数据库删除出错
	ErrorDatabaseDelete = 5013
	//数据库修改出错
	ErrorDatabaseUpdate = 5014

	//邮件服务相关错误
	ErrorEmailSend = 5020
)

// Response 是交付层的基本回应
type Response struct {
	Code    int         `json:"code"`    //请求状态代码
	Message interface{} `json:"message"` //请求结果提示
	Data    interface{} `json:"data"`    //请求结果与错误原因
}

// List 会返回给交付层一个列表回应
type List struct {
	Code    int         `json:"code"`    //请求状态代码
	Count   int         `json:"count"`   //数据量
	Message interface{} `json:"message"` //请求结果提示
	Data    interface{} `json:"data"`    //请求结果
}
