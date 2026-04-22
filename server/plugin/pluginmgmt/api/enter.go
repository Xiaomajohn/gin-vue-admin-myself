package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/service"
)

type ApiGroup struct {
	PluginApi
}

var ApiGroupApp = new(ApiGroup)

var pluginService = service.ServiceGroupApp.PluginService
