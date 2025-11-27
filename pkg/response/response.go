package response

import (
	"net/http"
	"project/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// Success 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// SuccessNoData 成功响应（无数据）
func SuccessNoData(c *gin.Context, code int, mess string) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: errcode.ErrorMessage[errcode.Success],
	})
}

// Failed 失败响应
func Failed(c *gin.Context, code int, message string) {
	c.JSON(errcode.Success, Response{
		Code:    code,
		Message: message,
	})
}

// FailedWithData 失败响应（带数据）
func FailedWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(errcode.Success, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    errcode.BadRequest,
		Message: errcode.ErrorMessage[errcode.BadRequest],
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    errcode.Unauthorized,
		Message: errcode.ErrorMessage[errcode.Unauthorized],
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    errcode.Forbidden,
		Message: errcode.ErrorMessage[errcode.Forbidden],
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    errcode.NotFound,
		Message: errcode.ErrorMessage[errcode.NotFound],
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    errcode.ServerError,
		Message: errcode.ErrorMessage[errcode.ServerError],
	})
}
