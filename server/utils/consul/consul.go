package consul

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var ConsulClient *api.Client

// InitConsul 初始化Consul客户端
func InitConsul() error {
	config := api.DefaultConfig()
	config.Address = global.GVA_CONFIG.Consul.Address
	if global.GVA_CONFIG.Consul.Token != "" {
		config.Token = global.GVA_CONFIG.Consul.Token
	}

	client, err := api.NewClient(config)
	if err != nil {
		global.GVA_LOG.Error("创建Consul客户端失败", zap.Error(err))
		return err
	}

	ConsulClient = client
	global.GVA_LOG.Info("Consul客户端初始化成功", zap.String("address", global.GVA_CONFIG.Consul.Address))
	return nil
}

// RegisterService 注册服务到Consul
func RegisterService(serviceName, serviceID, address string, port int, tags []string, healthCheckURL string) error {
	if ConsulClient == nil {
		return fmt.Errorf("Consul客户端未初始化")
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
	}

	// 配置健康检查
	if healthCheckURL != "" {
		registration.Checks = api.AgentServiceChecks{
			&api.AgentServiceCheck{
				HTTP:     healthCheckURL,
				Interval: "10s",
				Timeout:  "5s",
			},
		}
	}

	err := ConsulClient.Agent().ServiceRegister(registration)
	if err != nil {
		global.GVA_LOG.Error("注册服务到Consul失败",
			zap.String("service", serviceName),
			zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("服务注册成功", zap.String("service", serviceName))
	return nil
}

// DeregisterService 从Consul注销服务
func DeregisterService(serviceID string) error {
	if ConsulClient == nil {
		return fmt.Errorf("Consul客户端未初始化")
	}

	err := ConsulClient.Agent().ServiceDeregister(serviceID)
	if err != nil {
		global.GVA_LOG.Error("从Consul注销服务失败",
			zap.String("serviceID", serviceID),
			zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("服务注销成功", zap.String("serviceID", serviceID))
	return nil
}

// GetService 获取服务实例列表
func GetService(serviceName string) ([]*api.ServiceEntry, error) {
	if ConsulClient == nil {
		return nil, fmt.Errorf("Consul客户端未初始化")
	}

	entries, _, err := ConsulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		global.GVA_LOG.Error("获取服务实例失败",
			zap.String("service", serviceName),
			zap.Error(err))
		return nil, err
	}

	return entries, nil
}

// GetAllServices 获取所有服务
func GetAllServices() (map[string]*api.AgentService, error) {
	if ConsulClient == nil {
		return nil, fmt.Errorf("Consul客户端未初始化")
	}

	services, err := ConsulClient.Agent().Services()
	if err != nil {
		global.GVA_LOG.Error("获取所有服务失败", zap.Error(err))
		return nil, err
	}

	return services, nil
}

// CheckServiceHealth 检查服务健康状态
func CheckServiceHealth(serviceID string) (string, error) {
	if ConsulClient == nil {
		return "", fmt.Errorf("Consul客户端未初始化")
	}

	checks, _, err := ConsulClient.Health().Checks(serviceID, nil)
	if err != nil {
		return "", err
	}

	if len(checks) == 0 {
		return "unknown", nil
	}

	// 检查所有健康检查项
	for _, check := range checks {
		if check.Status == api.HealthCritical {
			return "critical", nil
		} else if check.Status == api.HealthWarning {
			return "warning", nil
		}
	}

	return "healthy", nil
}
