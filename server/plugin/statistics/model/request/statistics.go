package model

// ESSearchRequest ES搜索请求
type ESSearchRequest struct {
	Index     string `json:"index" form:"index"`         // 索引名称
	Query     string `json:"query" form:"query"`         // 查询DSL(JSON格式)
	Page      int    `json:"page" form:"page"`           // 页码
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 每页数量
	TimeRange string `json:"timeRange" form:"timeRange"` // 时间范围
}

// CPUStatisticsRequest CPU统计请求
type CPUStatisticsRequest struct {
	TimeRange string `json:"timeRange" form:"timeRange"` // 时间范围(如: now-1h, now-24h)
	Interval  string `json:"interval" form:"interval"`   // 聚合间隔(如: 1m, 5m, 1h)
	HostName  string `json:"hostName" form:"hostName"`   // 主机名(可选)
}

// MemoryStatisticsRequest 内存统计请求
type MemoryStatisticsRequest struct {
	TimeRange string `json:"timeRange" form:"timeRange"`
	HostName  string `json:"hostName" form:"hostName"`
}

// FileIntegrityRequest 文件完整性统计请求
type FileIntegrityRequest struct {
	TimeRange string `json:"timeRange" form:"timeRange"`
	FilePath  string `json:"filePath" form:"filePath"`   // 文件路径过滤
	EventType string `json:"eventType" form:"eventType"` // 事件类型(created/updated/deleted)
	HostName  string `json:"hostName" form:"hostName"`
}

// ProcessStatisticsRequest 进程统计请求
type ProcessStatisticsRequest struct {
	TimeRange   string `json:"timeRange" form:"timeRange"`
	ProcessName string `json:"processName" form:"processName"` // 进程名过滤
	HostName    string `json:"hostName" form:"hostName"`
	TopN        int    `json:"topN" form:"topN"` // Top N进程
}

// AuditEventRequest 审计事件统计请求
type AuditEventRequest struct {
	TimeRange string `json:"timeRange" form:"timeRange"`
	Category  string `json:"category" form:"category"` // 事件类别
	Action    string `json:"action" form:"action"`     // 动作
	UserName  string `json:"userName" form:"userName"` // 用户名
	HostName  string `json:"hostName" form:"hostName"`
}
