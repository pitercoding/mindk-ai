package models

import "time"

type DocumentEmbedding struct {
	ID        int       `json:"id"`
	ChunkID   int       `json:"chunk_id"`
	Embedding string    `json:"embedding"`
	CreatedAt time.Time `json:"created_at"`

	DocumentID   int    `json:"document_id"`
	DocumentName string `json:"document_name"`
	ChunkIndex   int    `json:"chunk_index"`
	Content      string `json:"content"`
}
