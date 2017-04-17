package module

import (
	"github.com/magic003/alice"
	"github.com/magic003/alice/example/persist"
)

// PersistModule is the module for persistent APIs.
type PersistModule struct {
	alice.BaseModule
	Table string `alice:"Table"`
}

// WebPageDao returns the WebPageDao.
func (m *PersistModule) WebPageDao() persist.WebPageDao {
	return persist.NewWebPageDao(m.Table)
}
