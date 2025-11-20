package redisdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"project/basic/config"
	"project/basic/logger"
)

var redisEntity *Redis

type Redis struct {
	name string

	isInit   bool
	addr     string
	network  string
	username string
	password string
	db       int

	redisDB *redis.Client
}

func newRedis() *Redis {
	return &Redis{
		name:     "redis",
		isInit:   false,
		db:       config.AppConfig.Db.Redis.DB,
		addr:     config.AppConfig.Db.Redis.Addr,
		network:  config.AppConfig.Db.Redis.Network,
		username: config.AppConfig.Db.Redis.Username,
		password: config.AppConfig.Db.Redis.Password,
	}
}

func GetRedis() *Redis {

	if redisEntity != nil {
		return redisEntity
	}

	logger.Sugar.Infof("\t[component] redis is initiating...")

	//判断配置文件是否加载成功
	if config.AppConfig == nil || config.AppConfig.Db.Redis == nil {
		logger.Sugar.Errorf("\t[component] redis config load failed")
		return nil
	}

	redisEntity = newRedis()
	return redisEntity
}

func (r *Redis) GetClient() *redis.Client {
	return r.redisDB
}

func (r *Redis) GetName() string {
	return r.name
}

func (r *Redis) InitComponent() bool {

	//判断是否初始化
	if r.isInit {
		logger.Sugar.Infof("\t[component] %s is inited", r.name)
		return true
	}

	options := &redis.Options{
		Addr:     r.addr,
		Network:  r.network,
		Username: r.username,
		Password: r.password,
		DB:       r.db,
	}

	r.redisDB = redis.NewClient(options)

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	ctx := context.Background()
	_, err := r.redisDB.Ping(ctx).Result()
	if err != nil {
		logger.Sugar.Infof("\t[component] %s init failed: %s", r.name, err)
		return false
	}

	r.isInit = true

	logger.Sugar.Infof("\t[component] %s init success", r.name)

	return true
}

func (r *Redis) IsInitialize() bool {
	return r.isInit
}
