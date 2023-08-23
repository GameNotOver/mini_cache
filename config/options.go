package config

import (
	"mini_cache/cache/rediscache"
)

type Configs struct {
	Redis map[string]rediscache.Options `yaml:"redis" mapstructure:"redis"`
}
