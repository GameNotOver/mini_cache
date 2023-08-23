package provider

import (
	"mini_cache/utils"
	"time"
)

type Option struct {
	ID   string `yaml:"id" mapstructure:"id"`
	Type string `yaml:"type" mapstructure:"type"`
	TTL  string `yaml:"ttl" mapstructure:"ttl"`
	// 以下为 redis cache 选项
	Prefix    string `yaml:"prefix" mapstructure:"prefix"`
	BatchSize int    `yaml:"batch_size" mapstructure:"batch_size"`
	// 以下为 memory cache 选项
	Size            int  `yaml:"size"`
	AttachTenantKey bool `yaml:"attach_tenant_key" mapstructure:"attach_tenant_key"`
	Compression     bool `yaml:"compression" mapstructure:"compression"`
}

type CacheOptions struct {
	Caches []Option `yaml:"caches"`
}

func (co *Option) GetTTL() time.Duration {
	d, _ := utils.ParseDuration(co.TTL)
	return d
}
