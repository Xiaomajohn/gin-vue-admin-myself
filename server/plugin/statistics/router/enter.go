package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/statistics/api"
)

type RouterGroup struct {
	StatisticsRouter
}

var RouterGroupApp = new(RouterGroup)

var statisticsApi = api.ApiGroupApp.StatisticsApi
