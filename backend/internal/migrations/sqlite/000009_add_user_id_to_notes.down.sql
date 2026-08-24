DROP INDEX IF EXISTS idx_notes_user_id;

ALTER TABLE notes DROP COLUMN user_id;
