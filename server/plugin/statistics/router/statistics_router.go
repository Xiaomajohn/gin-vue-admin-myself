package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type StatisticsRouter struct{}

func (s *StatisticsRouter) InitStatisticsRouter(Router *gin.RouterGroup) {
	statisticsRouter := Router.Group("statistics").Use(middleware.OperationRecord())
	statisticsRouterWithoutRecord := Router.Group("statistics")

	statisticsRouter.POST("customQuery", statisticsApi.CustomQuery) // 自定义查询

	statisticsRouterWithoutRecord.GET("systemOverview", statisticsApi.GetSystemOverview)    // 系统概览
	statisticsRouterWithoutRecord.GET("cpu", statisticsApi.GetCPUStatistics)                // CPU统计
	statisticsRouterWithoutRecord.GET("memory", statisticsApi.GetMemoryStatistics)          // 内存统计
	statisticsRouterWithoutRecord.GET("fileIntegrity", statisticsApi.GetFileIntegrityStats) // 文件完整性
	statisticsRouterWithoutRecord.GET("process", statisticsApi.GetProcessStatistics)        // 进程统计
	statisticsRouterWithoutRecord.GET("audit", statisticsApi.GetAuditEventStats)            // 审计事件
}
