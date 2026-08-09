CREATE TABLE chat_messages_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE
);

INSERT INTO chat_messages_new (
        id,
        session_id,
        role,
        content,
        created_at
    )
SELECT id,
    session_id,
    role,
    content,
    created_at
FROM chat_messages
WHERE session_id IS NOT NULL;
DROP TABLE chat_messages;
ALTER TABLE chat_messages_new
    RENAME TO chat_messages;