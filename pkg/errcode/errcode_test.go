package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{
			name: "系统级 - 成功",
			code: Success,
			want: "success",
		},
		{
			name: "系统级 - 未授权",
			code: Unauthorized,
			want: "未授权",
		},
		{
			name: "用户模块 - 用户不存在",
			code: UserNotFound,
			want: "用户不存在",
		},
		{
			name: "订单模块 - 订单已支付",
			code: OrderPaid,
			want: "订单已支付",
		},
		{
			name: "文件上传 - 文件过大",
			code: FileTooLarge,
			want: "文件过大",
		},
		{
			name: "未知错误码",
			code: 99999,
			want: "未知错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMessage(tt.code)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("使用默认消息", func(t *testing.T) {
		err := New(UserNotFound)
		assert.Equal(t, UserNotFound, err.Code)
		assert.Equal(t, "用户不存在", err.Message)
		assert.Equal(t, "用户不存在", err.Error())
	})

	t.Run("使用自定义消息", func(t *testing.T) {
		customMsg := "Custom: user not found in DB"
		err := New(UserNotFound, customMsg)
		assert.Equal(t, UserNotFound, err.Code)
		assert.Equal(t, customMsg, err.Message)
		assert.Equal(t, customMsg, err.Error())
	})

	t.Run("自定义空消息应使用默认", func(t *testing.T) {
		err := New(Forbidden, "")
		assert.Equal(t, Forbidden, err.Code)
		assert.Equal(t, "禁止访问", err.Message) // 不应被空字符串覆盖
	})
}

func TestError_Error(t *testing.T) {
	err := &Error{
		Code:    ServerError,
		Message: "数据库连接失败",
	}
	assert.Equal(t, "数据库连接失败", err.Error())
}
