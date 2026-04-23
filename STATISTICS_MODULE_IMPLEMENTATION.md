# 数据统计模块 - 完整实施方案

> 📅 创建时间: 2026-04-20  
> 🎯 模块目标: 从Elasticsearch统计数据，实现服务器资源监控、文件修改监控和进程运行监控

---

## 📋 目录

1. [模块架构](#模块架构)
2. [已完成的文件](#已完成的文件)
3. [待完成的核心功能](#待完成的核心功能)
4. [ES查询DSL示例](#es查询dsl示例)
5. [前端页面设计](#前端页面设计)
6. [配置说明](#配置说明)
7. [下一步操作](#下一步操作)

---

## 模块架构

```
数据统计模块 (statistics)
├── 后端 (server/plugin/statistics/)
│   ├── model/
│   │   └── request/
│   │       └── statistics.go          ✅ 已完成 - 请求模型
│   ├── service/
│   │   ├── enter.go                   ⏳ 待创建
│   │   └── statistics_service.go      ⏳ 待创建 - 核心统计逻辑
│   ├── api/
│   │   ├── enter.go                   ⏳ 待创建
│   │   └── statistics_api.go          ⏳ 待创建 - HTTP接口
│   ├── router/
│   │   ├── enter.go                   ⏳ 待创建
│   │   └── statistics_router.go       ⏳ 待创建 - 路由定义
│   └── plugin.go                      ⏳ 待创建 - 插件入口
│
├── 工具类 (server/utils/elasticsearch/)
│   └── client.go                      ✅ 已完成 - ES客户端封装
│
├── 配置 (server/config/)
│   ├── elasticsearch.go               ✅ 已完成 - ES配置结构
│   └── config.go                      ✅ 已更新 - 添加ES配置
│
└── 前端 (web/src/plugin/statistics/)
    ├── api/
    │   └── statistics.js              ⏳ 待创建 - API封装
    └── view/
        ├── dashboard.vue              ⏳ 待创建 - 监控仪表板
        ├── cpu.vue                    ⏳ 待创建 - CPU监控
        ├── memory.vue                 ⏳ 待创建 - 内存监控
        ├── file-integrity.vue         ⏳ 待创建 - 文件完整性
        └── process.vue                ⏳ 待创建 - 进程监控
```

---

## 已完成的文件

### 1. ES配置结构
**文件**: `server/config/elasticsearch.go`
```go
type Elasticsearch struct {
    Hosts    []string // ES服务器列表
    Username string   // 用户名
    Password string   // 密码
    Timeout  int      // 超时时间(秒)
}
```

### 2. ES客户端封装
**文件**: `server/utils/elasticsearch/client.go`

**核心功能**:
- ✅ `InitES()` - 初始化ES客户端
- ✅ `CheckConnection()` - 检查连接状态
- ✅ `QueryDSL()` - 执行DSL查询
- ✅ `AggregateQuery()` - 执行聚合查询
- ✅ `CountDocuments()` - 统计文档数量
- ✅ `SearchWithQuery()` - 带条件搜索

### 3. 请求模型
**文件**: `server/plugin/statistics/model/request/statistics.go`

**包含的请求结构**:
- ✅ `ESSearchRequest` - 通用ES搜索
- ✅ `CPUStatisticsRequest` - CPU统计
- ✅ `MemoryStatisticsRequest` - 内存统计
- ✅ `FileIntegrityRequest` - 文件完整性
- ✅ `ProcessStatisticsRequest` - 进程统计
- ✅ `AuditEventRequest` - 审计事件

---

## 待完成的核心功能

### Service层实现

需要创建 `server/plugin/statistics/service/statistics_service.go`

**核心方法**:

```go
type StatisticsService struct{}

// GetCPUStatistics 获取CPU统计数据
func (s *StatisticsService) GetCPUStatistics(req request.CPUStatisticsRequest) (map[string]interface{}, error) {
    // 构建ES聚合查询DSL
    // 查询 auditmetricbeat-metrics-* 索引
    // 按时间间隔聚合 system.cpu.total.pct
}

// GetMemoryStatistics 获取内存统计数据
func (s *StatisticsService) GetMemoryStatistics(req request.MemoryStatisticsRequest) (map[string]interface{}, error) {
    // 查询 system.memory.used.pct
    // 返回时间序列数据
}

// GetFileIntegrityStats 获取文件完整性统计
func (s *StatisticsService) GetFileIntegrityStats(req request.FileIntegrityRequest) (map[string]interface{}, error) {
    // 查询 auditmetricbeat-fileintegrity-* 索引
    // 统计文件变更事件
    // 按文件路径、事件类型分组
}

// GetProcessStatistics 获取进程统计
func (s *StatisticsService) GetProcessStatistics(req request.ProcessStatisticsRequest) (map[string]interface{}, error) {
    // 查询 system.process.*
    // Top N 进程(CPU/内存)
}

// GetAuditEventStats 获取审计事件统计
func (s *StatisticsService) GetAuditEventStats(req request.AuditEventRequest) (map[string]interface{}, error) {
    // 查询 auditmetricbeat-audit-* 索引
    // 按事件类别、动作、用户分组
}

// GetSystemOverview 获取系统概览
func (s *StatisticsService) GetSystemOverview() (map[string]interface{}, error) {
    // 综合统计数据
    // - 当前CPU/内存使用率
    // - 最近文件变更数
    // - 活跃进程数
    // - 审计事件数
}
```

### 关键ES查询DSL

#### 1. CPU使用率时间序列

```json
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "range": {
            "@timestamp": {
              "gte": "now-24h",
              "lte": "now"
            }
          }
        },
        {
          "term": {
            "metricset.name": "cpu"
          }
        }
      ]
    }
  },
  "aggs": {
    "cpu_over_time": {
      "date_histogram": {
        "field": "@timestamp",
        "interval": "5m"
      },
      "aggs": {
        "avg_cpu": {
          "avg": {
            "field": "system.cpu.total.pct"
          }
        }
      }
    }
  }
}
```

#### 2. 文件变更统计

```json
{
  "size": 0,
  "query": {
    "bool": {
      "must": [
        {
          "range": {
            "@timestamp": {
              "gte": "now-7d"
            }
          }
        }
      ]
    }
  },
  "aggs": {
    "events_by_type": {
      "terms": {
        "field": "event.type",
        "size": 10
      }
    },
    "files_changed": {
      "terms": {
        "field": "file.path",
        "size": 20
      }
    }
  }
}
```

#### 3. Top 10 CPU进程

```json
{
  "size": 0,
  "query": {
    "range": {
      "@timestamp": {
        "gte": "now-1h"
      }
    }
  },
  "aggs": {
    "top_processes": {
      "terms": {
        "field": "system.process.name",
        "size": 10,
        "order": {
          "avg_cpu": "desc"
        }
      },
      "aggs": {
        "avg_cpu": {
          "avg": {
            "field": "system.process.cpu.total.pct"
          }
        },
        "avg_memory": {
          "avg": {
            "field": "system.process.memory.rss.bytes"
          }
        }
      }
    }
  }
}
```

---

## 前端页面设计

### 1. 监控仪表板 (dashboard.vue)

**布局**:
```
┌─────────────────────────────────────────┐
│  系统概览卡片                             │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐  │
│  │ CPU  │ │内存  │ │文件  │ │进程  │  │
│  │ 45%  │ │ 62%  │ │ 128  │ │ 256  │  │
│  └──────┘ └──────┘ └──────┘ └──────┘  │
├─────────────────────────────────────────┤
│  CPU/内存使用趋势图 (ECharts)             │
│  ┌─────────────────────────────────┐   │
│  │     折线图 - 24小时趋势          │   │
│  └─────────────────────────────────┘   │
├─────────────────────────────────────────┤
│  最近文件变更    │   Top进程列表        │
│  ┌──────────┐   │   ┌──────────┐      │
│  │ /etc/... │   │   │ nginx    │      │
│  │ /var/... │   │   │ mysql    │      │
│  └──────────┘   │   └──────────┘      │
└─────────────────────────────────────────┘
```

### 2. CPU监控页面 (cpu.vue)

**功能**:
- 实时CPU使用率折线图
- 各核心使用率对比
- CPU负载趋势
- 时间范围选择(1h/6h/24h/7d)

### 3. 文件完整性页面 (file-integrity.vue)

**功能**:
- 文件变更事件列表
- 按文件路径搜索
- 事件类型筛选(created/updated/deleted)
- 变更详情查看

### 4. 进程监控页面 (process.vue)

**功能**:
- Top进程表格(CPU/内存)
- 进程搜索
- 进程历史趋势
- 异常进程告警

---

## 配置说明

### 后端配置 (config.yaml)

```yaml
elasticsearch:
  hosts:
    - "http://localhost:9200"
  username: "elastic"        # 可选
  password: "your_password"  # 可选
  timeout: 30                # 超时时间(秒)
```

### 索引配置

根据 FIELDS_DOCUMENTATION.md，AuditMetricbeat 使用以下索引:

| 索引模式 | 数据类型 | 用途 |
|---------|---------|------|
| `auditmetricbeat-audit-*` | 审计事件 | 安全审计、合规检查 |
| `auditmetricbeat-fileintegrity-*` | 文件变更 | 文件完整性监控 |
| `auditmetricbeat-metrics-*` | 系统指标 | CPU/内存/进程等 |

---

## 下一步操作

### 立即需要完成的任务

1. **安装ES依赖**
   ```bash
   cd server
   go get github.com/olivere/elastic/v7
   go mod tidy
   ```

2. **创建Service层**
   - 文件: `server/plugin/statistics/service/statistics_service.go`
   - 实现所有统计方法

3. **创建API层**
   - 文件: `server/plugin/statistics/api/statistics_api.go`
   - 添加Swagger注释

4. **创建Router层**
   - 文件: `server/plugin/statistics/router/statistics_router.go`

5. **创建插件入口**
   - 文件: `server/plugin/statistics/plugin.go`

6. **注册插件**
   - 更新: `server/plugin/register.go`

7. **创建前端页面**
   - API封装: `web/src/plugin/statistics/api/statistics.js`
   - 仪表板: `web/src/plugin/statistics/view/dashboard.vue`
   - 其他监控页面

### 测试清单

- [ ] ES连接测试
- [ ] CPU统计查询测试
- [ ] 内存统计查询测试
- [ ] 文件完整性查询测试
- [ ] 进程统计查询测试
- [ ] 前端页面渲染测试
- [ ] ECharts图表显示测试

---

## API接口设计

### 后端路由

```
GET  /statistics/systemOverview     - 系统概览
GET  /statistics/cpu                - CPU统计
GET  /statistics/memory             - 内存统计
GET  /statistics/fileIntegrity      - 文件完整性
GET  /statistics/process            - 进程统计
GET  /statistics/audit              - 审计事件
POST /statistics/customQuery        - 自定义ES查询
```

### 前端API调用示例

```javascript
import { getCPUStatistics } from '@/plugin/statistics/api/statistics'

const cpuData = await getCPUStatistics({
  timeRange: 'now-24h',
  interval: '5m',
  hostName: 'server-01'
})
```

---

## 技术栈

### 后端
- **ES客户端**: `github.com/olivere/elastic/v7`
- **聚合查询**: ES Aggregations API
- **数据转换**: Go map/json

### 前端
- **图表库**: ECharts 5.5.1
- **HTTP请求**: Axios (通过 @/utils/request)
- **UI组件**: Element Plus
- **时间处理**: @/utils/date

---

## 注意事项

### 1. 性能优化
- 使用ES聚合而非全量查询
- 设置合理的size限制
- 添加时间范围过滤
- 使用缓存减少重复查询

### 2. 错误处理
- ES连接失败时返回友好提示
- 查询超时处理
- 空数据状态展示

### 3. 安全性
- ES查询参数验证
- 防止注入攻击
- 限制查询时间范围

### 4. 数据一致性
- 时间格式统一使用ISO 8601
- 百分比字段统一使用0-100范围
- 字节单位统一转换

---

## 参考文档

- [FIELDS_DOCUMENTATION.md](./FIELDS_DOCUMENTATION.md) - ES字段完整梳理
- [Elasticsearch官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [olivere/elastic Go客户端](https://github.com/olivere/elastic)

---

**创建时间**: 2026-04-20  
**状态**: 🟡 部分完成 - 基础设施已就绪，待实现核心业务逻辑
