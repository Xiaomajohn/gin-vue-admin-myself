package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Plugin 插件管理 结构体
type Plugin struct {
	global.GVA_MODEL
	Name         string `json:"name" form:"name" gorm:"column:name;comment:插件名称;type:varchar(100);uniqueIndex"`                 // 插件名称
	Code         string `json:"code" form:"code" gorm:"column:code;comment:插件编码;type:varchar(50);uniqueIndex"`                  // 插件编码
	Version      string `json:"version" form:"version" gorm:"column:version;comment:插件版本;type:varchar(20)"`                     // 插件版本
	Type         int    `json:"type" form:"type" gorm:"column:type;comment:插件类型(1:内部Go插件 2:外部插件)"`                              // 插件类型
	Status       int    `json:"status" form:"status" gorm:"column:status;comment:状态(1:启用 2:禁用 3:异常)"`                           // 状态
	Description  string `json:"description" form:"description" gorm:"column:description;comment:插件描述;type:text"`                // 插件描述
	Author       string `json:"author" form:"author" gorm:"column:author;comment:作者;type:varchar(50)"`                          // 作者
	Icon         string `json:"icon" form:"icon" gorm:"column:icon;comment:图标;type:varchar(200)"`                               // 图标
	Config       string `json:"config" form:"config" gorm:"column:config;comment:插件配置(JSON格式);type:text"`                       // 插件配置
	HealthURL    string `json:"healthURL" form:"healthURL" gorm:"column:health_url;comment:健康检查URL;type:varchar(200)"`          // 健康检查URL(外部插件)
	ServiceName  string `json:"serviceName" form:"serviceName" gorm:"column:service_name;comment:Consul服务名称;type:varchar(100)"` // Consul服务名称(外部插件)
	HealthStatus int    `json:"healthStatus" form:"healthStatus" gorm:"column:health_status;comment:健康状态(1:健康 2:异常 3:未知)"`      // 健康状态
}

// TableName 插件 Plugin自定义表名
func (Plugin) TableName() string {
	return "test_plugin"
}
