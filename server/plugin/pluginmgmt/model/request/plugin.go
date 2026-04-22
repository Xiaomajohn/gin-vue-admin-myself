package request

import (
	crq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreatePluginRequest 创建插件请求
type CreatePluginRequest struct {
	Name        string `json:"name" binding:"required" form:"name"` // 插件名称
	Code        string `json:"code" binding:"required" form:"code"` // 插件编码
	Version     string `json:"version" form:"version"`              // 插件版本
	Type        int    `json:"type" binding:"required" form:"type"` // 插件类型
	Description string `json:"description" form:"description"`      // 插件描述
	Author      string `json:"author" form:"author"`                // 作者
	Icon        string `json:"icon" form:"icon"`                    // 图标
	Config      string `json:"config" form:"config"`                // 插件配置
	HealthURL   string `json:"healthURL" form:"healthURL"`          // 健康检查URL
	ServiceName string `json:"serviceName" form:"serviceName"`      // Consul服务名称
}

// UpdatePluginRequest 更新插件请求
type UpdatePluginRequest struct {
	ID          uint   `json:"ID" binding:"required" gorm:"primarykey"` // ID
	Name        string `json:"name" binding:"required" form:"name"`     // 插件名称
	Code        string `json:"code" binding:"required" form:"code"`     // 插件编码
	Version     string `json:"version" form:"version"`                  // 插件版本
	Type        int    `json:"type" binding:"required" form:"type"`     // 插件类型
	Status      int    `json:"status" form:"status"`                    // 状态
	Description string `json:"description" form:"description"`          // 插件描述
	Author      string `json:"author" form:"author"`                    // 作者
	Icon        string `json:"icon" form:"icon"`                        // 图标
	Config      string `json:"config" form:"config"`                    // 插件配置
	HealthURL   string `json:"healthURL" form:"healthURL"`              // 健康检查URL
	ServiceName string `json:"serviceName" form:"serviceName"`          // Consul服务名称
}

// PluginSearch 搜索插件请求
type PluginSearch struct {
	crq.PageInfo
	Name         string `json:"name" form:"name"`                 // 插件名称
	Code         string `json:"code" form:"code"`                 // 插件编码
	Type         int    `json:"type" form:"type"`                 // 插件类型
	Status       int    `json:"status" form:"status"`             // 状态
	HealthStatus int    `json:"healthStatus" form:"healthStatus"` // 健康状态
}

// UpdatePluginStatusRequest 更新插件状态请求
type UpdatePluginStatusRequest struct {
	ID     uint `json:"ID" binding:"required"`     // ID
	Status int  `json:"status" binding:"required"` // 状态
}

// IdsReq 批量删除请求
type IdsReq struct {
	Ids []uint `json:"ids" form:"ids"`
}

// GetById 根据ID获取请求
type GetById struct {
	ID uint `json:"ID" form:"ID"` // 主键ID
}
