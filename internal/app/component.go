package app

import (
	"fmt"
	"project/pkg/database"
	"project/pkg/logger"
	"project/pkg/queue"
)

// SetComponent 设置组件列表
func (d *DefaultApp) SetComponent(components []string) {
	d.components = components
	logger.Sugar.Infof("[component] components set: %v", components)
}

// LoadComponents 加载组件
func (d *DefaultApp) LoadComponents() error {
	if len(d.components) == 0 {
		logger.Sugar.Info("[component] no components to load")
		return nil
	}

	logger.Sugar.Infof("[component] loading %d components...", len(d.components))

	for _, comp := range d.components {
		if err := d.loadComponent(comp); err != nil {
			return fmt.Errorf("failed to load component '%s': %w", comp, err)
		}
	}

	logger.Sugar.Info("[component] all components loaded successfully")
	return nil
}

// CloseComponents 加载组件
func (d *DefaultApp) CloseComponents() error {
	if len(d.components) == 0 {
		logger.Sugar.Info("[component] no components to load")
		return nil
	}

	logger.Sugar.Infof("[component] closing %d components...", len(d.components))

	for _, comp := range d.components {
		if err := d.closeComponent(comp); err != nil {
			return fmt.Errorf("failed to close component '%s': %w", comp, err)
		}
	}

	logger.Sugar.Info("[component] all components closed successfully")
	return nil
}

// loadComponent 加载单个组件
func (d *DefaultApp) loadComponent(comp string) error {
	logger.Sugar.Infof("[component] loading component: %s", comp)

	switch comp {
	case "mysql":
		m := database.GetMysql()
		if m == nil {
			return fmt.Errorf("mysql component is nil")
		}
		m.InitComponent()

	case "gorm":
		gm := database.GetGormMysql()
		if gm == nil {
			return fmt.Errorf("gorm_mysql component is nil")
		}
		gm.InitComponent()

	case "redis":
		r := database.GetRedis()
		if r == nil {
			return fmt.Errorf("redis component is nil")
		}
		r.InitComponent()

	case "tdengine":
		t := database.GetTdengine()
		if t == nil {
			return fmt.Errorf("tdengine component is nil")
		}
		t.InitComponent()

	case "rabbitmq":
		rb := queue.GetRabbitMQ()
		if rb == nil {
			return fmt.Errorf("rabbitmq component is nil")
		}
		rb.InitComponent()

	default:
		return fmt.Errorf("unknown component: %s", comp)
	}

	logger.Sugar.Infof("[component] component '%s' loaded successfully", comp)
	return nil
}

// closeComponent 关闭单个组件
func (d *DefaultApp) closeComponent(comp string) error {
	logger.Sugar.Infof("[component] closing component: %s", comp)

	switch comp {
	case "mysql":
		m := database.GetMysql()
		if m == nil {
			return fmt.Errorf("mysql component is nil")
		}
		m.Close()

	case "gorm":
		gm := database.GetGormMysql()
		if gm == nil {
			return fmt.Errorf("gorm_mysql component is nil")
		}
		gm.Close()

	case "redis":
		r := database.GetRedis()
		if r == nil {
			return fmt.Errorf("redis component is nil")
		}
		r.Close()

	case "tdengine":
		t := database.GetTdengine()
		if t == nil {
			return fmt.Errorf("tdengine component is nil")
		}
		t.Close()

	case "rabbitmq":
		rb := queue.GetRabbitMQ()
		if rb == nil {
			return fmt.Errorf("rabbitmq component is nil")
		}
		rb.Close()

	default:
		return fmt.Errorf("unknown component: %s", comp)
	}

	logger.Sugar.Infof("[component] component '%s' closed successfully", comp)
	return nil
}
