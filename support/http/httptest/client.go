package httptest

// On is the entrypoint method into this packages "client" mocking system.
func (c *Client) On(method string, url string) *ClientExpectation {
	return &ClientExpectation{
		Method: method,
		URL:    url,
		Client: c,
	}
}

// OnWithBody is the entrypoint method into this packages "client" mocking system.
func (c *Client) OnWithBody(method string, url string, body []byte) *ClientExpectation {
	return &ClientExpectation{
		Method: method,
		URL:    url,
		Body:   body,
		Client: c,
	}
}
