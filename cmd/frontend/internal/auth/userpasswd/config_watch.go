package userpasswd

import (
	"github.com/sourcegraph/sourcegraph/cmd/frontend/auth/providers"
	"github.com/sourcegraph/sourcegraph/internal/conf"
)

// Watch for configuration changes related to the builtin auth provider.
func init() {
	go func() {
		conf.Watch(func() {
			newPC, _ := GetProviderConfig()
			if newPC == nil {
				providers.Update("builtin", nil)
				return
			}
			providers.Update("builtin", []providers.Provider{&provider{c: newPC}})
		})
	}()
}
