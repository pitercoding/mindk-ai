package services

type EmbeddingGenerator interface {
	Generate(text string) ([]float32, error)
}
