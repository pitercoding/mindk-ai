package llm

type OpenAIEmbeddingGenerator struct {
	client Client
}

func NewOpenAIEmbeddingGenerator(
	client Client,
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
