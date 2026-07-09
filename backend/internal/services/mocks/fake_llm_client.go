package mocks

type FakeLLMClient struct {
	Response   string
	Err        error
	LastPrompt string
}

func (c *FakeLLMClient) Chat(prompt string) (string, error) {
	c.LastPrompt = prompt
	return c.Response, c.Err
}
