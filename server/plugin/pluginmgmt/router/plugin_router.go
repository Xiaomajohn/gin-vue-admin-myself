package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PluginRouter struct{}

func (s *PluginRouter) InitPluginRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	pluginRouter := Router.Group("plugin").Use(middleware.OperationRecord())
	pluginRouterWithoutRecord := Router.Group("plugin")

	pluginRouter.POST("createPlugin", pluginApi.CreatePlugin)   // 新建插件
	pluginRouter.DELETE("deletePlugin", pluginApi.DeletePlugin) // 删除插件
	pluginRouter.PUT("updatePlugin", pluginApi.UpdatePlugin)    // 更新插件
	pluginRouter.PUT("updatePluginStatus", pluginApi.UpdatePluginStatus)

	pluginRouterWithoutRecord.GET("getPlugin", pluginApi.GetPlugin)         // 根据ID获取插件
	pluginRouterWithoutRecord.GET("getPluginList", pluginApi.GetPluginList) // 获取插件列表
}
