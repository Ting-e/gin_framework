package errcode

// 错误码定义
const (
	// 系统级错误 (0-999)
	Success        = 200
	BadRequest     = 400
	Unauthorized   = 401
	Forbidden      = 403
	NotFound       = 404
	RequestTimeout = 408
	TooManyRequest = 429
	ServerError    = 500

	// 用户模块 (1000-1999)（示例）
	UserNotFound      = 1001
	UserAlreadyExists = 1002
	UserPasswordError = 1003
	UserDisabled      = 1004

	// 订单模块 (2000-2999)（示例）
	OrderNotFound  = 2001
	OrderCancelled = 2002
	OrderPaid      = 2003

	// 文件上传 (4000-4999)（示例）
	FileUploadFailed   = 4001
	FileTooLarge       = 4002
	FileTypeNotAllowed = 4003
)

// ErrorMessage 错误码对应的消息
var ErrorMessage = map[int]string{
	// 系统级
	Success:        "success",
	BadRequest:     "请求参数错误",
	Unauthorized:   "未授权",
	Forbidden:      "禁止访问",
	NotFound:       "资源不存在",
	ServerError:    "服务器错误",
	RequestTimeout: "请求超时",
	TooManyRequest: "请求过多",

	// 用户模块（示例）
	UserNotFound:      "用户不存在",
	UserAlreadyExists: "用户已存在",
	UserPasswordError: "密码错误",
	UserDisabled:      "用户已被禁用",

	// 订单模块（示例）
	OrderNotFound:  "订单不存在",
	OrderCancelled: "订单已取消",
	OrderPaid:      "订单已支付",

	// 文件上传（示例）
	FileUploadFailed:   "文件上传失败",
	FileTooLarge:       "文件过大",
	FileTypeNotAllowed: "不支持的文件类型",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := ErrorMessage[code]; ok {
		return msg
	}
	return "未知错误"
}

// Error 自定义错误类型
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// New 创建错误
func New(code int, message ...string) *Error {
	msg := GetMessage(code)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}
