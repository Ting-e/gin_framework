package model

type PageData struct {
	Title         string // 页面标题
	ServiceName   string // 服务名称
	Now           string // 服务器时间
	Environment   string // 环境 (Development/Production)
	Version       string // 版本号
	StatusMessage string // 状态消息
}
