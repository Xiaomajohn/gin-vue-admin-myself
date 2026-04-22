package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PluginApi struct{}

// CreatePlugin 创建插件
// @Tags     PluginMgmt
// @Summary  创建插件
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreatePluginRequest true "插件信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /plugin/createPlugin [post]
func (p *PluginApi) CreatePlugin(c *gin.Context) {
	var req request.CreatePluginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	plugin := model.Plugin{
		Name:        req.Name,
		Code:        req.Code,
		Version:     req.Version,
		Type:        req.Type,
		Status:      1,
		Description: req.Description,
		Author:      req.Author,
		Icon:        req.Icon,
		Config:      req.Config,
		HealthURL:   req.HealthURL,
		ServiceName: req.ServiceName,
	}

	err = pluginService.CreatePlugin(plugin)
	if err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeletePlugin 删除插件
// @Tags     PluginMgmt
// @Summary  删除插件
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.GetById true "插件ID数组"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /plugin/deletePlugin [delete]
func (p *PluginApi) DeletePlugin(c *gin.Context) {
	var req request.IdsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = pluginService.DeletePlugin(req.Ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdatePlugin 更新插件
// @Tags     PluginMgmt
// @Summary  更新插件
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdatePluginRequest true "插件信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /plugin/updatePlugin [put]
func (p *PluginApi) UpdatePlugin(c *gin.Context) {
	var req request.UpdatePluginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	plugin := model.Plugin{
		Name:        req.Name,
		Code:        req.Code,
		Version:     req.Version,
		Type:        req.Type,
		Status:      req.Status,
		Description: req.Description,
		Author:      req.Author,
		Icon:        req.Icon,
		Config:      req.Config,
		HealthURL:   req.HealthURL,
		ServiceName: req.ServiceName,
	}
	plugin.ID = req.ID

	err = pluginService.UpdatePlugin(plugin)
	if err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetPlugin 获取插件详情
// @Tags     PluginMgmt
// @Summary  获取插件详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    ID query uint true "插件ID"
// @Success  200  {object} response.Response{data=model.Plugin,msg=string} "获取成功"
// @Router   /plugin/getPlugin [get]
func (p *PluginApi) GetPlugin(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	plugin, err := pluginService.GetPlugin(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(plugin, c)
}

// GetPluginList 获取插件列表
// @Tags     PluginMgmt
// @Summary  分页获取插件列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.PluginSearch true "查询参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /plugin/getPluginList [get]
func (p *PluginApi) GetPluginList(c *gin.Context) {
	var req request.PluginSearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := pluginService.GetPluginInfoList(req)
	if err != nil {
		global.GVA_LOG.Error("获取失败", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// UpdatePluginStatus 更新插件状态
// @Tags     PluginMgmt
// @Summary  更新插件状态
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdatePluginStatusRequest true "插件ID和状态"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /plugin/updatePluginStatus [put]
func (p *PluginApi) UpdatePluginStatus(c *gin.Context) {
	var req request.UpdatePluginStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = pluginService.UpdatePluginStatus(req.ID, req.Status)
	if err != nil {
		global.GVA_LOG.Error("更新状态失败", zap.Error(err))
		response.FailWithMessage("更新状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
