ALTER TABLE documents ADD COLUMN user_id TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_documents_user_id ON documents(user_id);
