package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/statistics/service"
)

type ApiGroup struct {
	StatisticsApi
}

var ApiGroupApp = new(ApiGroup)

var statisticsService = service.ServiceGroupApp.StatisticsService
