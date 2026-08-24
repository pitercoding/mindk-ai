ALTER TABLE chat_sessions ADD COLUMN user_id TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_chat_sessions_user_id ON chat_sessions(user_id);
