package init

import (
	"mini_cache/common"
	"mini_cache/di"
)

func Init() {
	di.MustRegister(common.NewProviderFromConfig) // cache provider
}
