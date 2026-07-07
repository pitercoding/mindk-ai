package llm

type Client interface {
    Chat(prompt string) (string, error)
}