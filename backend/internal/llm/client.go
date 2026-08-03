package llm

type Client interface {
	Chat(prompt string) (string, error)

	CreateEmbedding(text string) ([]float32, error)
}
