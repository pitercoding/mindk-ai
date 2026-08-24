DROP INDEX IF EXISTS idx_chat_sessions_user_id;

ALTER TABLE chat_sessions DROP COLUMN user_id;
