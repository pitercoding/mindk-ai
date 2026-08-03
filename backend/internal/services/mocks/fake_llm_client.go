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

func (f *FakeLLMClient) CreateEmbedding(
	text string,
) ([]float32, error) {

	return nil, nil
}
