CREATE TABLE document_embeddings (
    id SERIAL PRIMARY KEY,
    chunk_id INTEGER NOT NULL,
    embedding TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chunk_id) REFERENCES document_chunks(id) ON DELETE CASCADE
);

CREATE INDEX idx_document_embeddings_chunk_id
ON document_embeddings(chunk_id);