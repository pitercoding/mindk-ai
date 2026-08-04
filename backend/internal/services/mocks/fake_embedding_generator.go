package mocks

type FakeEmbeddingGenerator struct {
	Called bool
	Texts  []string
	Vector []float32
	Err    error
}

func (f *FakeEmbeddingGenerator) Generate(
	text string,
) ([]float32, error) {

	f.Called = true

	f.Texts = append(
		f.Texts,
		text,
	)

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Vector, nil
}
