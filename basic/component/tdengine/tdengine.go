package tdenginedb

import (
	"database/sql"
	"project/basic/config"
	"project/basic/logger"
)

var tdengineEntity *Tdengine

type Tdengine struct {
	//组件名称
	name string
	//clickhouse url
	url string
	//初始化标识
	isInit bool
	//conn
	taos *sql.DB
}

func newTdengine() *Tdengine {
	return &Tdengine{
		name:   "tdengine",
		isInit: false,
		url:    config.GetConfig().Db.Tdengine.Url,
	}
}

func GetTdengine() *Tdengine {

	if tdengineEntity != nil {
		return tdengineEntity
	}

	//判断配置文件是否加载
	if config.GetConfig() == nil || config.GetConfig().Minio == nil {
		logger.Sugar.Errorf("\t[component] clickhouseEntity config load failed")
		return nil
	}

	tdengineEntity = newTdengine()
	tdengineEntity.InitComponent()

	return tdengineEntity
}

func (t *Tdengine) GetName() string {
	return t.name
}

func (t *Tdengine) InitComponent() bool {

	//判断是否初始化
	if t.isInit {
		logger.Sugar.Infof("\t[component] %s is inited", t.name)
		return true
	}

	logger.Sugar.Infof("\t[component] %s is initiating...", t.name)

	var err error
	t.taos, err = sql.Open("taosRestful", t.url)
	if err != nil {
		logger.Sugar.Errorf("\t[component] %s init failed: %s", t.name, err)
		return false
	}
	// 激活连接
	if err = t.taos.Ping(); err != nil {
		logger.Sugar.Fatalf("\t[component] %s connect failed: %s", t.name, err)
		return false
	}
	//将初始化设置为true
	t.isInit = true

	logger.Sugar.Infof("\t[component] %s init success", t.name)

	return true
}

func (t *Tdengine) IsInitialize() bool {
	return t.isInit
}

func (t *Tdengine) GetUrl() string {
	return t.url
}

// 获取db
func (t *Tdengine) GetDb() *sql.DB {
	return t.taos
}
