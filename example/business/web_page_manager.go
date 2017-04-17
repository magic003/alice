package business

import (
	"github.com/magic003/alice/example/client"
	"github.com/magic003/alice/example/persist"
)

// WebPageManager manages web pages.
type WebPageManager struct {
	httpClient client.HTTPClient
	webPageDao persist.WebPageDao
}

// NewWebPageManager returns a new instance of WebPageManager.
func NewWebPageManager(httpClient client.HTTPClient, webPageDao persist.WebPageDao) *WebPageManager {
	if httpClient == nil || webPageDao == nil {
		panic("httpClient or webPageDao cannot be nil")
	}
	return &WebPageManager{
		httpClient: httpClient,
		webPageDao: webPageDao,
	}
}
