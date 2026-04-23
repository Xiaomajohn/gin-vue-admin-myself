package service

import (
	"encoding/json"
	"fmt"

	req "github.com/flipped-aurora/gin-vue-admin/server/plugin/statistics/model/request"
	es "github.com/flipped-aurora/gin-vue-admin/server/utils/elasticsearch"
	"github.com/olivere/elastic/v7"
)

type StatisticsService struct{}

// GetSystemOverview 获取系统概览
func (s *StatisticsService) GetSystemOverview() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取CPU使用率
	cpuQuery := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("@timestamp").Gte("now-5m")).
		Must(elastic.NewTermQuery("metricset.name", "cpu"))

	cpuCount, err := es.CountDocuments("auditmetricbeat-metrics-*", cpuQuery)
	if err != nil {
		return nil, err
	}
	result["cpuDataPoints"] = cpuCount

	// 获取文件变更数
	fileQuery := elastic.NewRangeQuery("@timestamp").Gte("now-1h")
	fileCount, err := es.CountDocuments("auditmetricbeat-fileintegrity-*", fileQuery)
	if err != nil {
		return nil, err
	}
	result["recentFileChanges"] = fileCount

	// 获取审计事件数
	auditQuery := elastic.NewRangeQuery("@timestamp").Gte("now-1h")
	auditCount, err := es.CountDocuments("auditmetricbeat-audit-*", auditQuery)
	if err != nil {
		return nil, err
	}
	result["recentAuditEvents"] = auditCount

	return result, nil
}

// GetCPUStatistics 获取CPU统计数据
func (s *StatisticsService) GetCPUStatistics(reqParams req.CPUStatisticsRequest) (map[string]interface{}, error) {
	// 构建聚合查询DSL
	aggQuery := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": reqParams.TimeRange,
								"lte": "now",
							},
						},
					},
					{
						"term": map[string]interface{}{
							"metricset.name": "cpu",
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"cpu_over_time": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":    "@timestamp",
					"interval": reqParams.Interval,
				},
				"aggs": map[string]interface{}{
					"avg_cpu": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "system.cpu.total.pct",
						},
					},
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(aggQuery)
	result, err := es.QueryDSL("auditmetricbeat-metrics-*", string(queryJSON))
	if err != nil {
		return nil, err
	}

	// 解析聚合结果
	response := make(map[string]interface{})
	if agg, found := result.Aggregations["cpu_over_time"]; found {
		response["timeSeries"] = agg
	}

	return response, nil
}

// GetMemoryStatistics 获取内存统计数据
func (s *StatisticsService) GetMemoryStatistics(reqParams req.MemoryStatisticsRequest) (map[string]interface{}, error) {
	aggQuery := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": reqParams.TimeRange,
								"lte": "now",
							},
						},
					},
					{
						"term": map[string]interface{}{
							"metricset.name": "memory",
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"memory_over_time": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":    "@timestamp",
					"interval": "5m",
				},
				"aggs": map[string]interface{}{
					"avg_memory_pct": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "system.memory.used.pct",
						},
					},
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(aggQuery)
	result, err := es.QueryDSL("auditmetricbeat-metrics-*", string(queryJSON))
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	if agg, found := result.Aggregations["memory_over_time"]; found {
		response["timeSeries"] = agg
	}

	return response, nil
}

// GetFileIntegrityStats 获取文件完整性统计
func (s *StatisticsService) GetFileIntegrityStats(reqParams req.FileIntegrityRequest) (map[string]interface{}, error) {
	boolQuery := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("@timestamp").Gte(reqParams.TimeRange))

	if reqParams.FilePath != "" {
		boolQuery.Must(elastic.NewWildcardQuery("file.path", fmt.Sprintf("*%s*", reqParams.FilePath)))
	}

	if reqParams.EventType != "" {
		boolQuery.Must(elastic.NewTermQuery("event.type", reqParams.EventType))
	}

	if reqParams.HostName != "" {
		boolQuery.Must(elastic.NewTermQuery("host.name", reqParams.HostName))
	}

	// 获取事件类型分布
	aggQuery := map[string]interface{}{
		"size":  0,
		"query": boolQuery,
		"aggs": map[string]interface{}{
			"events_by_type": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "event.type",
					"size":  10,
				},
			},
			"recent_changes": map[string]interface{}{
				"top_hits": map[string]interface{}{
					"sort": []map[string]interface{}{
						{"@timestamp": map[string]interface{}{"order": "desc"}},
					},
					"size": 20,
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(aggQuery)
	result, err := es.QueryDSL("auditmetricbeat-fileintegrity-*", string(queryJSON))
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	if agg, found := result.Aggregations["events_by_type"]; found {
		response["eventsByType"] = agg
	}
	if agg, found := result.Aggregations["recent_changes"]; found {
		response["recentChanges"] = agg
	}

	return response, nil
}

// GetProcessStatistics 获取进程统计
func (s *StatisticsService) GetProcessStatistics(reqParams req.ProcessStatisticsRequest) (map[string]interface{}, error) {
	if reqParams.TopN <= 0 {
		reqParams.TopN = 10
	}

	boolQuery := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("@timestamp").Gte(reqParams.TimeRange)).
		Must(elastic.NewTermQuery("metricset.name", "process"))

	if reqParams.ProcessName != "" {
		boolQuery.Must(elastic.NewWildcardQuery("system.process.name", fmt.Sprintf("*%s*", reqParams.ProcessName)))
	}

	if reqParams.HostName != "" {
		boolQuery.Must(elastic.NewTermQuery("host.name", reqParams.HostName))
	}

	aggQuery := map[string]interface{}{
		"size":  0,
		"query": boolQuery,
		"aggs": map[string]interface{}{
			"top_processes": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "system.process.name",
					"size":  reqParams.TopN,
					"order": map[string]interface{}{
						"avg_cpu": "desc",
					},
				},
				"aggs": map[string]interface{}{
					"avg_cpu": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "system.process.cpu.total.pct",
						},
					},
					"avg_memory": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "system.process.memory.rss.bytes",
						},
					},
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(aggQuery)
	result, err := es.QueryDSL("auditmetricbeat-metrics-*", string(queryJSON))
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	if agg, found := result.Aggregations["top_processes"]; found {
		response["topProcesses"] = agg
	}

	return response, nil
}

// GetAuditEventStats 获取审计事件统计
func (s *StatisticsService) GetAuditEventStats(reqParams req.AuditEventRequest) (map[string]interface{}, error) {
	boolQuery := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("@timestamp").Gte(reqParams.TimeRange))

	if reqParams.Category != "" {
		boolQuery.Must(elastic.NewTermQuery("event.category", reqParams.Category))
	}

	if reqParams.Action != "" {
		boolQuery.Must(elastic.NewTermQuery("event.action", reqParams.Action))
	}

	if reqParams.UserName != "" {
		boolQuery.Must(elastic.NewTermQuery("user.name", reqParams.UserName))
	}

	if reqParams.HostName != "" {
		boolQuery.Must(elastic.NewTermQuery("host.name", reqParams.HostName))
	}

	aggQuery := map[string]interface{}{
		"size":  0,
		"query": boolQuery,
		"aggs": map[string]interface{}{
			"events_by_category": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "event.category",
					"size":  20,
				},
			},
			"events_by_action": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "event.action",
					"size":  20,
				},
			},
			"events_over_time": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":    "@timestamp",
					"interval": "1h",
				},
			},
		},
	}

	queryJSON, _ := json.Marshal(aggQuery)
	result, err := es.QueryDSL("auditmetricbeat-audit-*", string(queryJSON))
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	if agg, found := result.Aggregations["events_by_category"]; found {
		response["eventsByCategory"] = agg
	}
	if agg, found := result.Aggregations["events_by_action"]; found {
		response["eventsByAction"] = agg
	}
	if agg, found := result.Aggregations["events_over_time"]; found {
		response["eventsOverTime"] = agg
	}

	return response, nil
}

// CustomQuery 自定义ES查询
func (s *StatisticsService) CustomQuery(reqParams req.ESSearchRequest) (map[string]interface{}, error) {
	result, err := es.QueryDSL(reqParams.Index, reqParams.Query)
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	response["total"] = result.TotalHits()
	response["hits"] = result.Hits.Hits

	return response, nil
}
