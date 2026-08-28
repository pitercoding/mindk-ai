package mocks

type FakeLLMClient struct {
	Response   string
	Err        error
	LastPrompt string
}

func (f *FakeLLMClient) Chat(prompt string) (string, error) {
	f.LastPrompt = prompt
	return f.Response, f.Err
}

func (f *FakeLLMClient) CreateEmbedding(
	text string,
) ([]float32, error) {

	return nil, nil
}
