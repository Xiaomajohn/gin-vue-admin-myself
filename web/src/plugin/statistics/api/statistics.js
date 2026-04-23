import service from '@/utils/request'

/**
 * 获取系统概览
 * @returns {Promise}
 */
export const getSystemOverview = () => {
  return service({
    url: '/statistics/systemOverview',
    method: 'get'
  })
}

/**
 * 获取CPU统计数据
 * @param {Object} params 查询参数
 * @param {string} params.timeRange 时间范围(如: now-1h, now-24h)
 * @param {string} params.interval 聚合间隔(如: 1m, 5m, 1h)
 * @param {string} params.hostName 主机名(可选)
 * @returns {Promise}
 */
export const getCPUStatistics = (params) => {
  return service({
    url: '/statistics/cpu',
    method: 'get',
    params: params
  })
}

/**
 * 获取内存统计数据
 * @param {Object} params 查询参数
 * @param {string} params.timeRange 时间范围
 * @param {string} params.hostName 主机名
 * @returns {Promise}
 */
export const getMemoryStatistics = (params) => {
  return service({
    url: '/statistics/memory',
    method: 'get',
    params: params
  })
}

/**
 * 获取文件完整性统计
 * @param {Object} params 查询参数
 * @param {string} params.timeRange 时间范围
 * @param {string} params.filePath 文件路径过滤
 * @param {string} params.eventType 事件类型(created/updated/deleted)
 * @param {string} params.hostName 主机名
 * @returns {Promise}
 */
export const getFileIntegrityStats = (params) => {
  return service({
    url: '/statistics/fileIntegrity',
    method: 'get',
    params: params
  })
}

/**
 * 获取进程统计数据
 * @param {Object} params 查询参数
 * @param {string} params.timeRange 时间范围
 * @param {string} params.processName 进程名过滤
 * @param {string} params.hostName 主机名
 * @param {number} params.topN Top N进程
 * @returns {Promise}
 */
export const getProcessStatistics = (params) => {
  return service({
    url: '/statistics/process',
    method: 'get',
    params: params
  })
}

/**
 * 获取审计事件统计
 * @param {Object} params 查询参数
 * @param {string} params.timeRange 时间范围
 * @param {string} params.category 事件类别
 * @param {string} params.action 动作
 * @param {string} params.userName 用户名
 * @param {string} params.hostName 主机名
 * @returns {Promise}
 */
export const getAuditEventStats = (params) => {
  return service({
    url: '/statistics/audit',
    method: 'get',
    params: params
  })
}

/**
 * 执行自定义ES查询
 * @param {Object} data 查询参数
 * @param {string} data.index 索引名称
 * @param {string} data.query 查询DSL(JSON格式)
 * @param {number} data.page 页码
 * @param {number} data.pageSize 每页数量
 * @returns {Promise}
 */
export const customQuery = (data) => {
  return service({
    url: '/statistics/customQuery',
    method: 'post',
    data: data
  })
}
