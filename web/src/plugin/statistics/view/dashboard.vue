<template>
  <div class="statistics-dashboard">
    <!-- 系统概览卡片 -->
    <el-row :gutter="20" class="overview-cards">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="card-content">
            <div class="card-icon cpu-icon">
              <el-icon :size="40"><Cpu /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ overview.cpuUsage }}%</div>
              <div class="card-label">CPU使用率</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="card-content">
            <div class="card-icon memory-icon">
              <el-icon :size="40"><Memory /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ overview.memoryUsage }}%</div>
              <div class="card-label">内存使用率</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="card-content">
            <div class="card-icon file-icon">
              <el-icon :size="40"><Document /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ overview.fileChanges }}</div>
              <div class="card-label">文件变更(1h)</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="card-content">
            <div class="card-icon audit-icon">
              <el-icon :size="40"><List /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ overview.auditEvents }}</div>
              <div class="card-label">审计事件(1h)</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- CPU和内存趋势图 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>CPU使用趋势</span>
              <el-radio-group v-model="cpuTimeRange" size="small" @change="loadCPUData">
                <el-radio-button label="now-1h">1小时</el-radio-button>
                <el-radio-button label="now-6h">6小时</el-radio-button>
                <el-radio-button label="now-24h">24小时</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="cpuChartRef" class="chart-container"></div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>内存使用趋势</span>
              <el-radio-group v-model="memoryTimeRange" size="small" @change="loadMemoryData">
                <el-radio-button label="now-1h">1小时</el-radio-button>
                <el-radio-button label="now-6h">6小时</el-radio-button>
                <el-radio-button label="now-24h">24小时</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="memoryChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 文件变更和Top进程 -->
    <el-row :gutter="20" class="tables-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近文件变更</span>
              <el-button type="primary" size="small" @click="viewAllFileChanges">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentFileChanges" style="width: 100%" max-height="400">
            <el-table-column prop="file.path" label="文件路径" show-overflow-tooltip />
            <el-table-column prop="event.type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag v-if="row['event.type'] === 'creation'" type="success">创建</el-tag>
                <el-tag v-else-if="row['event.type'] === 'change'" type="warning">修改</el-tag>
                <el-tag v-else-if="row['event.type'] === 'deletion'" type="danger">删除</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="@timestamp" label="时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row['@timestamp']) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>Top 10 进程(CPU)</span>
              <el-button type="primary" size="small" @click="viewAllProcesses">查看全部</el-button>
            </div>
          </template>
          <el-table :data="topProcesses" style="width: 100%" max-height="400">
            <el-table-column prop="key" label="进程名" />
            <el-table-column prop="avg_cpu.value" label="CPU%" width="100">
              <template #default="{ row }">
                {{ (row.avg_cpu?.value || 0).toFixed(2) }}%
              </template>
            </el-table-column>
            <el-table-column prop="avg_memory.value" label="内存" width="120">
              <template #default="{ row }">
                {{ formatBytes(row.avg_memory?.value || 0) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Cpu, Memory, Document, List } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { formatTimeToStr } from '@/utils/date'
import {
  getSystemOverview,
  getCPUStatistics,
  getMemoryStatistics,
  getFileIntegrityStats,
  getProcessStatistics
} from '@/plugin/statistics/api/statistics'

// 概览数据
const overview = ref({
  cpuUsage: 0,
  memoryUsage: 0,
  fileChanges: 0,
  auditEvents: 0
})

// 时间范围
const cpuTimeRange = ref('now-24h')
const memoryTimeRange = ref('now-24h')

// 图表引用
const cpuChartRef = ref(null)
const memoryChartRef = ref(null)
let cpuChart = null
let memoryChart = null

// 表格数据
const recentFileChanges = ref([])
const topProcesses = ref([])

// 加载系统概览
const loadOverview = async () => {
  try {
    const res = await getSystemOverview()
    if (res.code === 0) {
      overview.value = res.data
    }
  } catch (error) {
    console.error('加载系统概览失败:', error)
  }
}

// 加载CPU数据
const loadCPUData = async () => {
  try {
    const res = await getCPUStatistics({
      timeRange: cpuTimeRange.value,
      interval: '5m'
    })
    if (res.code === 0 && cpuChart) {
      updateCPUChart(res.data)
    }
  } catch (error) {
    console.error('加载CPU数据失败:', error)
  }
}

// 加载内存数据
const loadMemoryData = async () => {
  try {
    const res = await getMemoryStatistics({
      timeRange: memoryTimeRange.value
    })
    if (res.code === 0 && memoryChart) {
      updateMemoryChart(res.data)
    }
  } catch (error) {
    console.error('加载内存数据失败:', error)
  }
}

// 加载文件变更
const loadFileChanges = async () => {
  try {
    const res = await getFileIntegrityStats({
      timeRange: 'now-1h'
    })
    if (res.code === 0) {
      recentFileChanges.value = res.data.recentChanges?.hits?.hits || []
    }
  } catch (error) {
    console.error('加载文件变更失败:', error)
  }
}

// 加载Top进程
const loadTopProcesses = async () => {
  try {
    const res = await getProcessStatistics({
      timeRange: 'now-1h',
      topN: 10
    })
    if (res.code === 0) {
      topProcesses.value = res.data.topProcesses?.buckets || []
    }
  } catch (error) {
    console.error('加载Top进程失败:', error)
  }
}

// 更新CPU图表
const updateCPUChart = (data) => {
  const buckets = data.timeSeries?.buckets || []
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}<br/>{a}: {c}%'
    },
    xAxis: {
      type: 'category',
      data: buckets.map(b => formatTimeToStr(b.key, 'yyyy-MM-dd hh:mm:ss'))
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [{
      name: 'CPU使用率',
      type: 'line',
      data: buckets.map(b => (b.avg_cpu?.value || 0).toFixed(2)),
      smooth: true,
      itemStyle: {
        color: '#409EFF'
      },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(64, 158, 255, 0.5)' },
          { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
        ])
      }
    }]
  }
  cpuChart.setOption(option)
}

// 更新内存图表
const updateMemoryChart = (data) => {
  const buckets = data.timeSeries?.buckets || []
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}<br/>{a}: {c}%'
    },
    xAxis: {
      type: 'category',
      data: buckets.map(b => formatTimeToStr(b.key, 'yyyy-MM-dd hh:mm:ss'))
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [{
      name: '内存使用率',
      type: 'line',
      data: buckets.map(b => (b.avg_memory_pct?.value || 0).toFixed(2)),
      smooth: true,
      itemStyle: {
        color: '#67C23A'
      },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(103, 194, 58, 0.5)' },
          { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
        ])
      }
    }]
  }
  memoryChart.setOption(option)
}

// 格式化时间
const formatTime = (timestamp) => {
  return formatTimeToStr(timestamp, 'yyyy-MM-dd hh:mm:ss')
}

// 格式化字节
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

// 查看全部文件变更
const viewAllFileChanges = () => {
  // 跳转到文件完整性页面
  console.log('跳转到文件完整性页面')
}

// 查看所有进程
const viewAllProcesses = () => {
  // 跳转到进程监控页面
  console.log('跳转到进程监控页面')
}

// 初始化
onMounted(() => {
  // 初始化图表
  if (cpuChartRef.value) {
    cpuChart = echarts.init(cpuChartRef.value)
  }
  if (memoryChartRef.value) {
    memoryChart = echarts.init(memoryChartRef.value)
  }

  // 加载数据
  loadOverview()
  loadCPUData()
  loadMemoryData()
  loadFileChanges()
  loadTopProcesses()

  // 定时刷新
  const timer = setInterval(() => {
    loadOverview()
    loadCPUData()
    loadMemoryData()
  }, 30000) // 30秒刷新一次

  onUnmounted(() => {
    clearInterval(timer)
    cpuChart?.dispose()
    memoryChart?.dispose()
  })
})
</script>

<style scoped>
.statistics-dashboard {
  padding: 20px;
}

.overview-cards {
  margin-bottom: 20px;
}

.card-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.card-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.cpu-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.memory-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.file-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.audit-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.card-info {
  flex: 1;
}

.card-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.card-label {
  font-size: 14px;
  color: #909399;
}

.charts-row,
.tables-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 350px;
  width: 100%;
}
</style>
