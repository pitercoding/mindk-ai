package models

import "time"

type DocumentEmbedding struct {
	ID        int       `json:"id"`
	ChunkID   int       `json:"chunk_id"`
	Embedding string    `json:"embedding"`
	CreatedAt time.Time `json:"created_at"`
}
