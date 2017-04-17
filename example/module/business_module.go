package module

import (
	"github.com/magic003/alice"
	"github.com/magic003/alice/example/business"
	"github.com/magic003/alice/example/client"
	"github.com/magic003/alice/example/persist"
)

// BusinessModule is the module for business objects.
type BusinessModule struct {
	alice.BaseModule
	WebPageDao persist.WebPageDao `alice:""`
	HTTPClient client.HTTPClient  `alice:"HTTPClient"`
}

// WebPageManager returns an instance of WebPageManager.
func (m *BusinessModule) WebPageManager() *business.WebPageManager {
	return business.NewWebPageManager(m.HTTPClient, m.WebPageDao)
}
