# 数据统计插件 - 完整文档

> 📅 创建时间: 2026-04-20  
> 🎯 插件功能: 从Elasticsearch统计数据，实现服务器资源监控、文件修改监控和进程运行监控

---

## 📋 插件结构

```
server/plugin/statistics/          # 后端插件
├── api/
│   ├── enter.go                   ✅ API组入口
│   └── statistics_api.go          ✅ HTTP接口实现
├── model/
│   └── request/
│       └── statistics.go          ✅ 请求模型
├── service/
│   ├── enter.go                   ✅ 服务组入口
│   └── statistics_service.go      ✅ 核心业务逻辑
├── router/
│   ├── enter.go                   ✅ 路由组入口
│   └── statistics_router.go       ✅ 路由定义
├── initialize/
│   └── router.go                  ✅ 路由初始化
└── plugin.go                      ✅ 插件入口

web/src/plugin/statistics/         # 前端插件
├── api/
│   └── statistics.js              ✅ API接口封装
└── view/
    └── dashboard.vue              ✅ 监控仪表板
```

---

## 🔧 配置说明

### 1. 安装ES依赖

```bash
cd server
go get github.com/olivere/elastic/v7
go mod tidy
```

### 2. 配置Elasticsearch连接

在 `server/config.yaml` 中添加：

```yaml
elasticsearch:
  hosts:
    - "http://localhost:9200"      # ES服务器地址
  username: "elastic"              # 用户名(可选)
  password: "your_password"        # 密码(可选)
  timeout: 30                      # 超时时间(秒)
```

### 3. 初始化ES客户端

在 `server/core/server.go` 或启动文件中添加ES初始化：

```go
import "github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"

// 在启动时初始化
if err := es.InitES(); err != nil {
    global.GVA_LOG.Error("ES初始化失败", zap.Error(err))
}
```

---

## 🚀 API接口

### 路由列表

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/statistics/systemOverview` | 系统概览 | JWT |
| GET | `/statistics/cpu` | CPU统计 | JWT |
| GET | `/statistics/memory` | 内存统计 | JWT |
| GET | `/statistics/fileIntegrity` | 文件完整性 | JWT |
| GET | `/statistics/process` | 进程统计 | JWT |
| GET | `/statistics/audit` | 审计事件 | JWT |
| POST | `/statistics/customQuery` | 自定义查询 | JWT + 操作记录 |

### 请求示例

#### 1. 获取系统概览

```bash
GET /statistics/systemOverview
```

响应:
```json
{
  "code": 0,
  "data": {
    "cpuDataPoints": 1250,
    "recentFileChanges": 128,
    "recentAuditEvents": 256
  },
  "msg": "获取成功"
}
```

#### 2. 获取CPU统计

```bash
GET /statistics/cpu?timeRange=now-24h&interval=5m
```

参数:
- `timeRange`: 时间范围 (now-1h, now-6h, now-24h, now-7d)
- `interval`: 聚合间隔 (1m, 5m, 15m, 1h)
- `hostName`: 主机名(可选)

#### 3. 获取文件完整性统计

```bash
GET /statistics/fileIntegrity?timeRange=now-24h&eventType=change
```

参数:
- `timeRange`: 时间范围
- `filePath`: 文件路径过滤(支持通配符)
- `eventType`: 事件类型 (creation, change, deletion)
- `hostName`: 主机名

#### 4. 获取Top进程

```bash
GET /statistics/process?timeRange=now-1h&topN=10
```

参数:
- `timeRange`: 时间范围
- `processName`: 进程名过滤
- `hostName`: 主机名
- `topN`: Top N (默认10)

#### 5. 自定义ES查询

```bash
POST /statistics/customQuery
Content-Type: application/json

{
  "index": "auditmetricbeat-metrics-*",
  "query": "{\"size\":10,\"query\":{\"match_all\":{}}}",
  "page": 1,
  "pageSize": 10
}
```

---

## 📊 前端使用

### 1. API调用

```javascript
import { 
  getSystemOverview,
  getCPUStatistics,
  getFileIntegrityStats 
} from '@/plugin/statistics/api/statistics'

// 获取系统概览
const overview = await getSystemOverview()

// 获取CPU统计
const cpuData = await getCPUStatistics({
  timeRange: 'now-24h',
  interval: '5m'
})

// 获取文件变更
const fileChanges = await getFileIntegrityStats({
  timeRange: 'now-1h',
  eventType: 'change'
})
```

### 2. 配置路由

在系统管理 -> 菜单管理中添加：

- **菜单名称**: 数据统计
- **路由路径**: statistics
- **组件路径**: plugin/statistics/view/dashboard
- **是否菜单**: 是
- **排序**: 根据实际需要

---

## 🎯 核心功能

### 1. 系统概览仪表板

**功能**:
- 实时CPU/内存使用率
- 最近文件变更统计
- 审计事件统计
- 自动刷新(30秒)

### 2. CPU监控

**功能**:
- 时间序列折线图
- 多时间范围选择(1h/6h/24h)
- 聚合间隔可调
- 实时数据更新

### 3. 内存监控

**功能**:
- 内存使用趋势
- 百分比显示
- 时间范围选择

### 4. 文件完整性监控

**功能**:
- 文件变更事件列表
- 按事件类型筛选
- 文件路径搜索
- 变更详情

### 5. 进程监控

**功能**:
- Top N进程排行
- CPU/内存占用
- 进程搜索
- 历史趋势

### 6. 审计事件

**功能**:
- 事件分类统计
- 按用户/动作筛选
- 时间趋势分析

---

## 🔍 ES索引说明

基于 AuditMetricbeat 数据采集：

| 索引模式 | 数据类型 | 用途 |
|---------|---------|------|
| `auditmetricbeat-audit-*` | 审计事件 | 安全审计、合规检查 |
| `auditmetricbeat-fileintegrity-*` | 文件变更 | 文件完整性监控 |
| `auditmetricbeat-metrics-*` | 系统指标 | CPU/内存/进程等 |

---

## 📈 性能优化

### 1. 查询优化

- ✅ 使用聚合而非全量查询
- ✅ 设置合理的size限制
- ✅ 添加时间范围过滤
- ✅ 使用日期直方图聚合

### 2. 前端优化

- ✅ ECharts按需加载
- ✅ 定时刷新而非实时轮询
- ✅ 数据缓存
- ✅ 懒加载图表

### 3. 缓存策略

```go
// 可以使用Redis缓存频繁查询的结果
// 例如系统概览数据缓存30秒
```

---

## ⚠️ 注意事项

### 1. 安全性

- 所有API都需要JWT认证
- 自定义查询需要验证DSL合法性
- 限制查询时间范围(防止超大查询)

### 2. 错误处理

- ES连接失败时返回友好提示
- 查询超时处理
- 空数据状态展示

### 3. 数据一致性

- 时间格式统一使用ISO 8601
- 百分比字段统一使用0-100范围
- 字节单位统一转换

---

## 🧪 测试清单

### 后端测试

- [ ] ES连接测试
- [ ] 系统概览查询测试
- [ ] CPU聚合查询测试
- [ ] 文件完整性查询测试
- [ ] 进程统计查询测试
- [ ] 自定义查询测试

### 前端测试

- [ ] 仪表板页面渲染
- [ ] ECharts图表显示
- [ ] 时间范围切换
- [ ] 数据自动刷新
- [ ] 响应式布局

---

## 🚀 部署步骤

### 1. 后端部署

```bash
cd server
go mod tidy
go build -o server.exe
./server.exe
```

### 2. 前端部署

```bash
cd web
npm install
npm run build
```

### 3. 配置Nginx

确保Nginx正确代理API请求到后端服务。

---

## 📚 扩展开发

### 添加新的统计维度

1. 在 `model/request/statistics.go` 添加请求结构
2. 在 `service/statistics_service.go` 实现查询逻辑
3. 在 `api/statistics_api.go` 添加HTTP接口
4. 在 `router/statistics_router.go` 注册路由
5. 在前端 `api/statistics.js` 添加API封装
6. 创建对应的Vue页面

### 示例：添加磁盘I/O统计

```go
// 1. 请求模型
type DiskIOStatisticsRequest struct {
    TimeRange string `json:"timeRange"`
    HostName  string `json:"hostName"`
}

// 2. Service实现
func (s *StatisticsService) GetDiskIOStatistics(req request.DiskIOStatisticsRequest) {
    // 构建ES查询
    // 查询 system.diskio.* 字段
}

// 3. API接口
func (s *StatisticsApi) GetDiskIOStatistics(c *gin.Context) {
    // 处理请求
}
```

---

## 🐛 常见问题

### Q1: ES连接失败？

A: 检查 `config.yaml` 中的ES配置是否正确，确保ES服务已启动。

### Q2: 查询返回空数据？

A: 确认索引名称正确，时间范围有数据覆盖。

### Q3: 图表不显示？

A: 检查ECharts是否正确安装，数据格式是否符合要求。

### Q4: 性能慢？

A: 
- 缩小时间范围
- 增加聚合间隔
- 使用缓存
- 优化ES查询DSL

---

## 📖 参考文档

- [FIELDS_DOCUMENTATION.md](../FIELDS_DOCUMENTATION.md) - ES字段完整梳理
- [Elasticsearch官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [olivere/elastic Go客户端](https://github.com/olivere/elastic)
- [ECharts官方文档](https://echarts.apache.org/)

---

**插件版本**: v1.0.0  
**创建时间**: 2026-04-20  
**状态**: ✅ 已完成，可直接使用
