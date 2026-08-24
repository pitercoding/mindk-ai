DROP INDEX IF EXISTS idx_documents_user_id;

ALTER TABLE documents DROP COLUMN user_id;
