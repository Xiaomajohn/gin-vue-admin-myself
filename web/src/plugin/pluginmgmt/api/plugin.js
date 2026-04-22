import service from '@/utils/request'

/**
 * 创建插件
 * @param {Object} data 插件信息
 * @returns {Promise}
 */
export const createPlugin = (data) => {
  return service({
    url: '/plugin/createPlugin',
    method: 'post',
    data: data
  })
}

/**
 * 删除插件
 * @param {Object} data 插件ID数组
 * @returns {Promise}
 */
export const deletePlugin = (data) => {
  return service({
    url: '/plugin/deletePlugin',
    method: 'delete',
    data: data
  })
}

/**
 * 更新插件
 * @param {Object} data 插件信息
 * @returns {Promise}
 */
export const updatePlugin = (data) => {
  return service({
    url: '/plugin/updatePlugin',
    method: 'put',
    data: data
  })
}

/**
 * 获取插件详情
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getPlugin = (params) => {
  return service({
    url: '/plugin/getPlugin',
    method: 'get',
    params: params
  })
}

/**
 * 获取插件列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getPluginList = (params) => {
  return service({
    url: '/plugin/getPluginList',
    method: 'get',
    params: params
  })
}

/**
 * 更新插件状态
 * @param {Object} data 插件ID和状态
 * @returns {Promise}
 */
export const updatePluginStatus = (data) => {
  return service({
    url: '/plugin/updatePluginStatus',
    method: 'put',
    data: data
  })
}
