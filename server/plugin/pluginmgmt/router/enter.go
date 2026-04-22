package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/api"
)

type RouterGroup struct {
	PluginRouter
}

var RouterGroupApp = new(RouterGroup)

var pluginApi = api.ApiGroupApp.PluginApi
