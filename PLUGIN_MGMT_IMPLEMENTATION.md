# 测试管理平台 - 插件管理系统实施文档

## 📋 项目概述

基于 gin-vue-admin 框架开发的测试管理平台插件管理系统，支持内部 Go 插件和外部多语言插件的统一管理。

## 🎯 核心功能

### Phase 1: 基础插件管理（已完成 ✅）

#### 1. 数据模型设计

**插件表 (test_plugin)**
- 支持内部Go插件和外部插件两种类型
- 完整的插件信息管理（名称、编码、版本、作者、描述等）
- 状态管理（启用/禁用/异常）
- 健康状态监控
- Consul服务集成配置

**文件位置：**
- 数据模型：`server/plugin/pluginmgmt/model/plugin.go`
- 请求模型：`server/plugin/pluginmgmt/model/request/plugin.go`

#### 2. 后端实现

**Service 层** (`server/plugin/pluginmgmt/service/`)
- ✅ CreatePlugin - 创建插件
- ✅ DeletePlugin - 删除插件（支持批量）
- ✅ UpdatePlugin - 更新插件
- ✅ GetPlugin - 获取插件详情
- ✅ GetPluginInfoList - 分页查询插件列表
- ✅ UpdatePluginStatus - 更新插件状态
- ✅ CheckPluginCodeUnique - 检查插件编码唯一性

**API 层** (`server/plugin/pluginmgmt/api/`)
- ✅ 完整的 Swagger 文档注释
- ✅ 统一错误处理
- ✅ 参数验证

**Router 层** (`server/plugin/pluginmgmt/router/`)
- ✅ RESTful API 路由设计
- ✅ JWT 认证中间件
- ✅ Casbin 权限控制
- ✅ 操作记录中间件

**路由列表：**
```
POST   /plugin/createPlugin       - 创建插件
DELETE /plugin/deletePlugin       - 删除插件
PUT    /plugin/updatePlugin       - 更新插件
PUT    /plugin/updatePluginStatus - 更新插件状态
GET    /plugin/getPlugin          - 获取插件详情
GET    /plugin/getPluginList      - 获取插件列表
```

#### 3. 前端实现

**API 封装** (`web/src/plugin/pluginmgmt/api/plugin.js`)
- ✅ 统一的 HTTP 请求封装
- ✅ 完整的 JSDoc 注释

**管理页面** (`web/src/plugin/pluginmgmt/view/plugin.vue`)
- ✅ 搜索/过滤功能
- ✅ 数据表格展示
- ✅ 新增/编辑插件（Drawer）
- ✅ 批量删除
- ✅ 状态切换
- ✅ 分页功能
- ✅ 响应式设计

### Phase 2: Consul 集成（已完成 ✅）

#### 1. Consul 客户端封装

**文件位置：** `server/utils/consul/consul.go`

**核心功能：**
- ✅ InitConsul - 初始化 Consul 客户端
- ✅ RegisterService - 注册服务到 Consul
- ✅ DeregisterService - 从 Consul 注销服务
- ✅ GetService - 获取服务实例列表
- ✅ GetAllServices - 获取所有服务
- ✅ CheckServiceHealth - 检查服务健康状态

#### 2. 配置支持

**配置文件：** `server/config/consul.go`
```yaml
consul:
  address: "localhost:8500"
  token: ""  # 可选
```

**主配置集成：** `server/config/config.go`
- ✅ 添加 Consul 配置字段

## 📁 项目结构

```
server/plugin/pluginmgmt/
├── api/
│   ├── enter.go              # API组入口
│   └── plugin_api.go         # 插件API实现
├── initialize/
│   ├── gorm.go              # 数据库初始化
│   └── router.go            # 路由初始化
├── model/
│   ├── plugin.go            # 数据模型
│   └── request/
│       └── plugin.go        # 请求模型
├── router/
│   ├── enter.go             # 路由组入口
│   └── plugin_router.go     # 路由定义
├── service/
│   ├── enter.go             # 服务组入口
│   └── plugin_service.go    # 服务实现
└── plugin.go                # 插件入口

server/utils/consul/
└── consul.go                # Consul客户端封装

server/config/
├── consul.go                # Consul配置
└── config.go                # 主配置（已更新）

web/src/plugin/pluginmgmt/
├── api/
│   └── plugin.js            # API接口封装
└── view/
    └── plugin.vue           # 插件管理页面
```

## 🚀 部署指南

### 1. 后端部署

#### 1.1 安装 Consul 依赖

```bash
cd server
go get github.com/hashicorp/consul/api
```

#### 1.2 配置 Consul

在 `server/config.yaml` 中添加：

```yaml
consul:
  address: "localhost:8500"  # Consul 服务器地址
  token: ""                   # Consul Token（可选）
```

#### 1.3 启动 Consul（本地开发）

```bash
# 使用 Docker 启动 Consul
docker run -d --name consul -p 8500:8500 consul:latest agent -server -bootstrap -ui -client=0.0.0.0
```

访问 Consul UI: http://localhost:8500

#### 1.4 编译并启动服务

```bash
cd server
go mod tidy
go run main.go
```

### 2. 前端部署

```bash
cd web
npm install
npm run dev
```

### 3. 配置菜单和权限

登录系统后，在"系统管理" -> "菜单管理"中添加插件管理菜单：

- **菜单名称**: 插件管理
- **路由路径**: plugin
- **组件路径**: plugin/pluginmgmt/view/plugin
- **排序**: 根据实际需要设置

## 🔧 使用示例

### 1. 创建内部 Go 插件

内部插件已在编译时自动注册，示例：

```go
// server/plugin/pluginmgmt/plugin.go
package pluginmgmt

func init() {
    interfaces.Register(Plugin)
}
```

### 2. 创建外部插件（Python 示例）

#### 2.1 创建 Python 插件服务

```python
from flask import Flask, jsonify

app = Flask(__name__)

@app.route('/health')
def health():
    return jsonify({"status": "healthy"})

@app.route('/api/execute', methods=['POST'])
def execute():
    # 执行测试逻辑
    return jsonify({"result": "success"})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

#### 2.2 注册到 Consul

使用插件管理界面添加：
- **插件名称**: Python测试插件
- **插件编码**: python-test-plugin
- **插件类型**: 外部插件 (2)
- **健康检查URL**: http://localhost:5000/health
- **Consul服务名**: python-test-plugin

### 3. 使用 Consul 进行服务发现

```go
import "github.com/flipped-aurora/gin-vue-admin/server/utils/consul"

// 获取插件服务实例
services, err := consul.GetService("python-test-plugin")
if err != nil {
    // 处理错误
}

// 调用插件服务
for _, service := range services {
    url := fmt.Sprintf("http://%s:%d/api/execute", 
        service.Service.Address, 
        service.Service.Port)
    // 发起 HTTP 请求
}
```

## 🎨 扩展建议

### Phase 3: 高级功能（未来规划）

1. **插件执行引擎**
   - 统一的插件调用接口
   - 异步任务队列
   - 执行日志记录

2. **插件市场**
   - 插件上传和下载
   - 版本管理
   - 评分和评论

3. **监控和告警**
   - 插件性能监控
   - 异常告警
   - 使用统计

4. **插件 SDK**
   - Python SDK
   - Java SDK
   - Node.js SDK
   - 开发文档

5. **安全性增强**
   - 插件沙箱
   - 权限细粒度控制
   - 审计日志

## ⚠️ 注意事项

### 1. 数据类型一致性

前后端数据类型必须保持一致：
- Go 的 `uint` 对应 JavaScript 的 `number`
- 状态字段使用整数枚举
- 时间字段使用标准格式

### 2. 错误处理

- Service 层返回 error
- API 层统一处理并返回格式化响应
- 前端统一错误提示

### 3. 权限控制

- 所有 API 都需要 JWT 认证
- 使用 Casbin 进行权限控制
- 按钮级别的权限控制

### 4. 健康检查

- 外部插件必须提供健康检查端点
- Consul 定期检查服务健康状态
- 健康状态同步到数据库

## 📊 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: Gin 1.10.0
- **ORM**: GORM 1.25.12
- **服务发现**: HashiCorp Consul
- **权限**: Casbin 2.103.0

### 前端
- **框架**: Vue 3.5.7
- **构建工具**: Vite 6.2.3
- **UI组件**: Element Plus 2.10.2
- **状态管理**: Pinia 2.2.2
- **HTTP客户端**: Axios 1.8.2

## 🎓 学习资源

- [Gin-Vue-Admin 官方文档](https://www.gin-vue-admin.com/)
- [Consul 官方文档](https://www.consul.io/docs)
- [GORM 官方文档](https://gorm.io/)
- [Vue 3 官方文档](https://cn.vuejs.org/)

## 📝 更新日志

### v1.0.0 (2025-04-20)
- ✅ 完成插件管理基础功能
- ✅ 实现 CRUD 操作
- ✅ 集成 Consul 服务发现
- ✅ 支持健康检查
- ✅ 完整的前后端实现

## 💡 常见问题

### Q1: Consul 连接失败？
A: 检查 config.yaml 中的 consul.address 配置是否正确，确保 Consul 服务已启动。

### Q2: 插件表未创建？
A: 确保数据库已初始化，插件会在启动时自动创建表结构。

### Q3: 前端页面 404？
A: 检查菜单配置中的组件路径是否正确：`plugin/pluginmgmt/view/plugin`

### Q4: 权限验证失败？
A: 确保用户有对应的 API 权限，在"系统管理" -> "API管理"中配置。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
