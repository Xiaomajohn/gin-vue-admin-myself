package config

type Consul struct {
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Token   string `mapstructure:"token" json:"token" yaml:"token"`
}
