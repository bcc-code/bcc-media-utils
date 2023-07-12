package cantemo

import "net/http"

type ItemsClient struct {
	*Client
}

func (c *Client) Items() *ItemsClient {
	return &ItemsClient{c}
}

func (c *ItemsClient) Get(id string) (*Item, error) {
	return request[Item](c.Client, http.MethodGet, c.url("API", "v2", "items", id).String(), nil)
}

func (c *ItemsClient) GetMetadata(id string) (*ItemMetadata, error) {
	return request[ItemMetadata](c.Client, http.MethodGet, c.url("API", "v2", "items", id, "metadata").String(), nil)
}
