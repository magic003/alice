package persist

// WebPageDao defines the interface to access web pages saved in database.
type WebPageDao interface {
	// Find looks up a web page by url.
	Find(url string) string
	// Save saves a web page.
	Save(url string, content string) error
}

// NewWebPageDao returns an instance of WebPageDao.
func NewWebPageDao(table string) WebPageDao {
	if table == "" {
		panic("table cannot be empty")
	}
	return &webPageDao{
		table: table,
	}
}

// webPageDao is a dummy implementation for WebPageDao.
type webPageDao struct {
	table string
}

func (dao *webPageDao) Find(url string) string {
	return ""
}

func (dao *webPageDao) Save(url string, content string) error {
	return nil
}
