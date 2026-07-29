package models

import "time"

type DocumentChunk struct {
	ID         int       `json:"id"`
	DocumentID int       `json:"document_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
