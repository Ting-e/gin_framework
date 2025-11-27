package queue

import (
	"project/pkg/config"
	"project/pkg/logger"

	"github.com/streadway/amqp"
)

var rabbitmqEntity *RabbitMQ

type RabbitMQ struct {
	//组件名称
	name string
	//clickhouse url
	url string
	//初始化标识
	isInit bool
	//conn
	rmq *amqp.Connection
}

func newRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		name:   "rabbitmq",
		isInit: false,
		url:    config.Get().RabbitMQ.URL,
	}
}

func GetRabbitMQ() *RabbitMQ {

	if rabbitmqEntity != nil {
		return rabbitmqEntity
	}

	//判断配置文件是否加载
	if config.Get() == nil || config.Get().RabbitMQ == nil {
		logger.Sugar.Errorf("\t[component] rabbitmqEntity config load failed")
		return nil
	}

	rabbitmqEntity = newRabbitMQ()

	return rabbitmqEntity
}

func (r *RabbitMQ) GetName() string {
	return r.name
}

func (r *RabbitMQ) InitComponent() bool {

	//判断是否初始化
	if r.isInit {
		logger.Sugar.Infof("\t[component] %s is inited", r.name)
		return true
	}

	logger.Sugar.Infof("\t[component] %s is initiating...", r.name)

	var err error
	r.rmq, err = amqp.Dial(r.url)
	if err != nil {
		logger.Sugar.Errorf("\t[component] %s init failed: %s", r.name, err)
		return false
	}
	// defer conn.Close()
	//将初始化设置为true
	r.isInit = true

	logger.Sugar.Infof("\t[component] %s init success", r.name)

	return true
}

// Close 关闭 RabbitMQ 连接
func (r *RabbitMQ) Close() error {
	if r.rmq == nil {
		return nil
	}
	err := r.rmq.Close()
	if err != nil {
		logger.Sugar.Errorf("\t[component] %s close failed: %v", r.name, err)
		return err
	}
	r.isInit = false
	r.rmq = nil // 避免重复关闭
	logger.Sugar.Infof("\t[component] %s closed successfully", r.name)
	return nil
}

func (r *RabbitMQ) IsInitialize() bool {
	return r.isInit
}

func (r *RabbitMQ) GetUrl() string {
	return r.url
}

// 获取db
func (r *RabbitMQ) GetQueue() *amqp.Connection {
	return r.rmq
}
