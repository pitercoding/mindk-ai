CREATE TABLE chat_messages_old (
    id SERIAL PRIMARY KEY,
    note_id INTEGER NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);

DROP TABLE chat_messages;

ALTER TABLE chat_messages_old
RENAME TO chat_messages;