ALTER TABLE chat_messages
ADD COLUMN session_id INTEGER REFERENCES chat_sessions(id) ON DELETE CASCADE;