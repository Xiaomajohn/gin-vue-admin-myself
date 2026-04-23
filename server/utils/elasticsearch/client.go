package es

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

var ESClient *elastic.Client

// InitES 初始化Elasticsearch客户端
func InitES() error {
	if global.GVA_CONFIG.Elasticsearch.Hosts == nil || len(global.GVA_CONFIG.Elasticsearch.Hosts) == 0 {
		global.GVA_LOG.Warn("Elasticsearch配置未设置，跳过初始化")
		return nil
	}

	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(global.GVA_CONFIG.Elasticsearch.Hosts...))
	options = append(options, elastic.SetSniff(false))
	options = append(options, elastic.SetHealthcheckInterval(time.Second*time.Duration(global.GVA_CONFIG.Elasticsearch.Timeout)))

	if global.GVA_CONFIG.Elasticsearch.Username != "" {
		options = append(options, elastic.SetBasicAuth(global.GVA_CONFIG.Elasticsearch.Username, global.GVA_CONFIG.Elasticsearch.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		global.GVA_LOG.Error("创建Elasticsearch客户端失败", zap.Error(err))
		return err
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	info, code, err := client.Ping(global.GVA_CONFIG.Elasticsearch.Hosts[0]).Do(ctx)
	if err != nil {
		global.GVA_LOG.Error("Elasticsearch连接测试失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info("Elasticsearch连接成功",
		zap.String("版本", info.Version.Number),
		zap.Int("状态码", code))

	ESClient = client
	return nil
}

// CheckConnection 检查ES连接状态
func CheckConnection() bool {
	if ESClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, _, err := ESClient.Ping(global.GVA_CONFIG.Elasticsearch.Hosts[0]).Do(ctx)
	return err == nil
}

// QueryDSL 执行DSL查询
func QueryDSL(index string, body string) (*elastic.SearchResult, error) {
	if ESClient == nil {
		return nil, fmt.Errorf("Elasticsearch客户端未初始化")
	}

	result, err := ESClient.Search(index).
		Source(body).
		Do(context.Background())

	if err != nil {
		global.GVA_LOG.Error("ES查询失败",
			zap.String("index", index),
			zap.Error(err))
		return nil, err
	}

	return result, nil
}

// AggregateQuery 执行聚合查询
func AggregateQuery(index string, body string) (map[string]interface{}, error) {
	result, err := QueryDSL(index, body)
	if err != nil {
		return nil, err
	}

	aggregations := make(map[string]interface{})
	for name, agg := range result.Aggregations {
		aggregations[name] = agg
	}

	return aggregations, nil
}

// CountDocuments 统计文档数量
func CountDocuments(index string, query elastic.Query) (int64, error) {
	if ESClient == nil {
		return 0, fmt.Errorf("Elasticsearch客户端未初始化")
	}

	count, err := ESClient.Count(index).Query(query).Do(context.Background())
	if err != nil {
		global.GVA_LOG.Error("ES统计失败",
			zap.String("index", index),
			zap.Error(err))
		return 0, err
	}

	return count, nil
}

// SearchWithQuery 带查询条件的搜索
func SearchWithQuery(index string, query elastic.Query, page, size int) (*elastic.SearchResult, error) {
	if ESClient == nil {
		return nil, fmt.Errorf("Elasticsearch客户端未初始化")
	}

	from := (page - 1) * size

	result, err := ESClient.Search(index).
		Query(query).
		From(from).
		Size(size).
		Sort("@timestamp", false).
		Do(context.Background())

	if err != nil {
		global.GVA_LOG.Error("ES搜索失败",
			zap.String("index", index),
			zap.Error(err))
		return nil, err
	}

	return result, nil
}
