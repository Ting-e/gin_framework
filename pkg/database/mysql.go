package database

import (
	"database/sql"
	"project/pkg/config"
	"project/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlEntity *Mysql

type Mysql struct {
	//组件名称
	name string
	//mysql url
	url string
	//初始化标识
	isInit bool
	//最大闲置数
	maxIdleConnection int
	//最大连接数
	maxOpenConnection int
	//db
	dB *sql.DB
}

func newMysql() *Mysql {
	return &Mysql{
		name:              "mysql",
		isInit:            false,
		maxOpenConnection: config.Get().Db.Mysql.MaxOpenConnection,
		maxIdleConnection: config.Get().Db.Mysql.MaxIdleConnection,
		url:               config.Get().Db.Mysql.URL,
	}
}

func GetMysql() *Mysql {

	if mysqlEntity != nil {
		return mysqlEntity
	}

	//判断配置文件是否加载
	if config.Get() == nil || config.Get().Db.Mysql == nil {
		logger.Sugar.Errorf("\t[component] mysql config load failed")
		return nil
	}

	mysqlEntity = newMysql()
	return mysqlEntity
}

func (m *Mysql) GetName() string {
	return m.name
}

func (m *Mysql) InitComponent() bool {

	//判断是否初始化
	if m.isInit {
		logger.Sugar.Infof("\t[component] %s is inited", m.name)
		return true
	}

	var err error
	m.dB, err = sql.Open("mysql", m.url)
	if err != nil {
		logger.Sugar.Errorf("\t[component] %s init failed: %s", m.name, err)
		return false
	}

	// 最大连接数
	m.dB.SetMaxOpenConns(m.maxOpenConnection)
	// 最大闲置数
	m.dB.SetMaxIdleConns(m.maxIdleConnection)

	// 激活连接
	if err = m.dB.Ping(); err != nil {
		logger.Sugar.Fatalf("\t[component] %s connect failed: %s", m.name, err)
		return false
	}

	//将初始化设置为true
	m.isInit = true

	return true
}

func (m *Mysql) IsInitialize() bool {
	return m.isInit
}

func (m *Mysql) GetMaxOpenConnection() int {
	return m.maxOpenConnection
}

func (m *Mysql) GetUrl() string {
	return m.url
}

// 获取db
func (m *Mysql) GetDb() *sql.DB {
	return m.dB
}
