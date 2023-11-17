package cantemo

import (
	"fmt"
	"net/http"
	"strings"
)

type SearchClient struct {
	*Client
}

func (c *Client) Search() *SearchClient {
	return &SearchClient{c}
}

func (c *SearchClient) Put(query string, page int) (*SearchResult, error) {
	u := c.url("API", "v2", "search/")
	u.RawQuery = fmt.Sprintf("page=%d", page)

	return request[SearchResult](c.Client, http.MethodPut, u.String(), strings.NewReader(query))
}
