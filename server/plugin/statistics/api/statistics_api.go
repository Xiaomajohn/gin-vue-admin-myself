package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	req "github.com/flipped-aurora/gin-vue-admin/server/plugin/statistics/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatisticsApi struct{}

// GetSystemOverview 获取系统概览
// @Tags     Statistics
// @Summary  获取系统概览数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/systemOverview [get]
func (s *StatisticsApi) GetSystemOverview(c *gin.Context) {
	data, err := statisticsService.GetSystemOverview()
	if err != nil {
		global.GVA_LOG.Error("获取系统概览失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetCPUStatistics 获取CPU统计
// @Tags     Statistics
// @Summary  获取CPU统计数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query req.CPUStatisticsRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/cpu [get]
func (s *StatisticsApi) GetCPUStatistics(c *gin.Context) {
	var requestParams req.CPUStatisticsRequest
	err := c.ShouldBindQuery(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if requestParams.TimeRange == "" {
		requestParams.TimeRange = "now-24h"
	}
	if requestParams.Interval == "" {
		requestParams.Interval = "5m"
	}

	data, err := statisticsService.GetCPUStatistics(requestParams)
	if err != nil {
		global.GVA_LOG.Error("获取CPU统计失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetMemoryStatistics 获取内存统计
// @Tags     Statistics
// @Summary  获取内存统计数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query req.MemoryStatisticsRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/memory [get]
func (s *StatisticsApi) GetMemoryStatistics(c *gin.Context) {
	var requestParams req.MemoryStatisticsRequest
	err := c.ShouldBindQuery(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if requestParams.TimeRange == "" {
		requestParams.TimeRange = "now-24h"
	}

	data, err := statisticsService.GetMemoryStatistics(requestParams)
	if err != nil {
		global.GVA_LOG.Error("获取内存统计失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetFileIntegrityStats 获取文件完整性统计
// @Tags     Statistics
// @Summary  获取文件完整性统计
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query req.FileIntegrityRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/fileIntegrity [get]
func (s *StatisticsApi) GetFileIntegrityStats(c *gin.Context) {
	var requestParams req.FileIntegrityRequest
	err := c.ShouldBindQuery(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if requestParams.TimeRange == "" {
		requestParams.TimeRange = "now-24h"
	}

	data, err := statisticsService.GetFileIntegrityStats(requestParams)
	if err != nil {
		global.GVA_LOG.Error("获取文件完整性统计失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetProcessStatistics 获取进程统计
// @Tags     Statistics
// @Summary  获取进程统计数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query req.ProcessStatisticsRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/process [get]
func (s *StatisticsApi) GetProcessStatistics(c *gin.Context) {
	var requestParams req.ProcessStatisticsRequest
	err := c.ShouldBindQuery(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if requestParams.TimeRange == "" {
		requestParams.TimeRange = "now-1h"
	}
	if requestParams.TopN <= 0 {
		requestParams.TopN = 10
	}

	data, err := statisticsService.GetProcessStatistics(requestParams)
	if err != nil {
		global.GVA_LOG.Error("获取进程统计失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetAuditEventStats 获取审计事件统计
// @Tags     Statistics
// @Summary  获取审计事件统计
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query req.AuditEventRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router   /statistics/audit [get]
func (s *StatisticsApi) GetAuditEventStats(c *gin.Context) {
	var requestParams req.AuditEventRequest
	err := c.ShouldBindQuery(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if requestParams.TimeRange == "" {
		requestParams.TimeRange = "now-24h"
	}

	data, err := statisticsService.GetAuditEventStats(requestParams)
	if err != nil {
		global.GVA_LOG.Error("获取审计事件统计失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// CustomQuery 自定义ES查询
// @Tags     Statistics
// @Summary  执行自定义ES查询
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body req.ESSearchRequest true "查询参数"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "查询成功"
// @Router   /statistics/customQuery [post]
func (s *StatisticsApi) CustomQuery(c *gin.Context) {
	var requestParams req.ESSearchRequest
	err := c.ShouldBindJSON(&requestParams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := statisticsService.CustomQuery(requestParams)
	if err != nil {
		global.GVA_LOG.Error("自定义查询失败", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "查询成功", c)
}
