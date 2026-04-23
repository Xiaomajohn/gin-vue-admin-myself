package config

// Elasticsearch Elasticsearch配置
type Elasticsearch struct {
	Hosts    []string `mapstructure:"hosts" json:"hosts" yaml:"hosts"`          // ES服务器列表
	Username string   `mapstructure:"username" json:"username" yaml:"username"` // 用户名
	Password string   `mapstructure:"password" json:"password" yaml:"password"` // 密码
	Timeout  int      `mapstructure:"timeout" json:"timeout" yaml:"timeout"`    // 超时时间(秒)
}
