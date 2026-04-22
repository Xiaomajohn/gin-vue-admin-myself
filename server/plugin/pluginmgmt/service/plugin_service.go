package service

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/pluginmgmt/model/request"
	"gorm.io/gorm"
)

type PluginService struct{}

// CreatePlugin 创建插件
func (s *PluginService) CreatePlugin(plugin model.Plugin) (err error) {
	if !errors.Is(global.GVA_DB.Where("code = ?", plugin.Code).First(&model.Plugin{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同插件编码")
	}
	err = global.GVA_DB.Create(&plugin).Error
	return err
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(IDs []uint) (err error) {
	err = global.GVA_DB.Where("id in ?", IDs).Delete(&model.Plugin{}).Error
	return err
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(plugin model.Plugin) (err error) {
	err = global.GVA_DB.Where("id = ?", plugin.ID).Save(&plugin).Error
	return err
}

// GetPlugin 根据ID获取插件
func (s *PluginService) GetPlugin(ID uint) (plugin model.Plugin, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&plugin).Error
	return
}

// GetPluginInfoList 分页获取插件列表
func (s *PluginService) GetPluginInfoList(info request.PluginSearch) (list []model.Plugin, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&model.Plugin{})
	var plugins []model.Plugin

	// 条件过滤
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Code != "" {
		db = db.Where("code LIKE ?", "%"+info.Code+"%")
	}
	if info.Type > 0 {
		db = db.Where("type = ?", info.Type)
	}
	if info.Status > 0 {
		db = db.Where("status = ?", info.Status)
	}
	if info.HealthStatus > 0 {
		db = db.Where("health_status = ?", info.HealthStatus)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&plugins).Error
	return plugins, total, err
}

// UpdatePluginStatus 更新插件状态
func (s *PluginService) UpdatePluginStatus(ID uint, status int) (err error) {
	err = global.GVA_DB.Model(&model.Plugin{}).Where("id = ?", ID).Update("status", status).Error
	return err
}

// CheckPluginCodeUnique 检查插件编码是否唯一
func (s *PluginService) CheckPluginCodeUnique(code string, ID uint) bool {
	var plugin model.Plugin
	err := global.GVA_DB.Where("code = ? AND id != ?", code, ID).First(&plugin).Error
	return errors.Is(err, gorm.ErrRecordNotFound)
}
