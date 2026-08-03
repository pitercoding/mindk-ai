package llm

type OpenAIEmbeddingGenerator struct {
	client *OpenAIClient
}

func NewOpenAIEmbeddingGenerator(
	client *OpenAIClient,
) *OpenAIEmbeddingGenerator {

	return &OpenAIEmbeddingGenerator{
		client: client,
	}
}

func (g *OpenAIEmbeddingGenerator) Generate(
	text string,
) ([]float32, error) {

	return g.client.CreateEmbedding(text)
}
