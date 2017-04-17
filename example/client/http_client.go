package client

// HTTPClient is an interface for fetching web page.
type HTTPClient interface {
	// Fetch fetches a web page by url.
	Fetch(url string) string
}

// NewHTTPClient returns an instance of HTTPClient.
func NewHTTPClient(retries int) HTTPClient {
	if retries <= 0 {
		panic("retries should be greater than 0")
	}
	return &httpClient{
		retries: retries,
	}
}

// httpClient is a dummy implementation of HTTPClient.
type httpClient struct {
	retries int
}

func (client *httpClient) Fetch(url string) string {
	return ""
}
