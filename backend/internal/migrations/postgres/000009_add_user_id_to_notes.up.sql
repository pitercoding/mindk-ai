ALTER TABLE notes ADD COLUMN user_id TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_notes_user_id ON notes(user_id);
