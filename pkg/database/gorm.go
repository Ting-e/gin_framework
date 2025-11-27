package database

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"project/pkg/config"
	local_logger "project/pkg/logger"
)

var gormMysqlEntity *GormMysql

type GormMysql struct {
	name              string
	url               string
	isInit            bool
	maxIdleConnection int
	maxOpenConnection int
	db                *gorm.DB // 注意
}

func newGormMysql() *GormMysql {
	cfg := config.Get().Db.Mysql
	return &GormMysql{
		name:              "gorm_mysql",
		url:               cfg.URL,
		maxIdleConnection: cfg.MaxIdleConnection,
		maxOpenConnection: cfg.MaxOpenConnection,
	}
}

func GetGormMysql() *GormMysql {
	if gormMysqlEntity != nil {
		return gormMysqlEntity
	}

	if config.Get() == nil || config.Get().Db.Mysql == nil {
		local_logger.Sugar.Errorf("\t[component] gorm_mysql config load failed")
		return nil
	}

	gormMysqlEntity = newGormMysql()
	return gormMysqlEntity
}

func (g *GormMysql) GetName() string {
	return g.name
}

func (g *GormMysql) InitComponent() bool {
	if g.isInit {
		local_logger.Sugar.Infof("\t[component] %s is inited", g.name)
		return true
	}

	// 构建 GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用表名复数化：gorm默认会把结构体名称+s得到的结果映射为数据库中的表名，设置该项后则不加s
		},
	}

	var err error
	g.db, err = gorm.Open(mysql.Open(g.url), gormConfig)
	if err != nil {
		local_logger.Sugar.Errorf("\t[component] %s init failed: %v", g.name, err)
		return false
	}

	// 配置底层 sql.DB 连接池
	sqlDB, err := g.db.DB()
	if err != nil {
		local_logger.Sugar.Errorf("\t[component] %s get raw db failed: %v", g.name, err)
		return false
	}

	sqlDB.SetMaxOpenConns(g.maxOpenConnection)
	sqlDB.SetMaxIdleConns(g.maxIdleConnection)
	sqlDB.SetConnMaxLifetime(time.Hour) // 可根据需要调整

	// 测试连接
	if err = sqlDB.Ping(); err != nil {
		local_logger.Sugar.Errorf("\t[component] %s connect failed: %v", g.name, err)
		return false
	}

	g.isInit = true
	return true
}

// 关闭 GORM MySQL 底层数据库连接
func (g *GormMysql) Close() error {
	if g.db == nil {
		return nil
	}
	sqlDB, err := g.db.DB()
	if err != nil {
		local_logger.Sugar.Errorf("\t[component] %s get raw db failed on close: %v", g.name, err)
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		local_logger.Sugar.Errorf("\t[component] %s close failed: %v", g.name, err)
		return err
	}
	g.isInit = false
	return nil
}

func (g *GormMysql) IsInitialize() bool {
	return g.isInit
}

func (g *GormMysql) GetDb() *gorm.DB {
	return g.db
}
