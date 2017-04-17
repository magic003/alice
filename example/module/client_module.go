package module

import (
	"github.com/magic003/alice"
	"github.com/magic003/alice/example/client"
)

// ClientModule is the module for clients.
type ClientModule struct {
	alice.BaseModule
	Retries int `alice:"Retries"`
}

// HTTPClient returns an instance fo HTTPClient.
func (m *ClientModule) HTTPClient() client.HTTPClient {
	return client.NewHTTPClient(m.Retries)
}
