package PO

const (
	Code_Succeed    = 1000
	Code_Fail       = -10001
	Code_DBFail     = -10002
	Code_AccessFail = -10003
	Code_NotLogin   = -10004
	Code_WxFail     = -20001
)

const (
	Data_Succeed = "Succeed!"
)

var (
	Succeed_Succeed = ServiceErr{Code_Succeed, "成功"}
)

var (
	Error_WrongParameter = ServiceErr{Code_Fail, "参数错误"}
	Error_FromModel      = ServiceErr{Code_DBFail, "数据库操作错误"}
	Error_CannotAccess   = ServiceErr{Code_AccessFail, "无权操作"}
	Error_WxOperation    = ServiceErr{Code_WxFail, "微信信息获取失败"}
	Error_NeedLogin      = ServiceErr{Code_NotLogin, "未登录"}
	Error_LoginFail      = ServiceErr{Code_NotLogin, "登录失败"}
)
