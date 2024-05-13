#!/bin/bash

# SQLite database file path
DATABASE="chat.db"

# Check if the database file already exists
if [ ! -f "$DATABASE" ]; then
    # Create SQLite database file
    touch "$DATABASE"
    echo "SQLite database '$DATABASE' created."
else
    echo "SQLite database '$DATABASE' already exists."
    exit 0
fi

# SQLite commands to create table
sqlite3 "$DATABASE" <<EOF
CREATE TABLE IF NOT EXISTS conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message BLOB,
    conversationId TEXT
);
EOF

echo "Table 'conversations' created in '$DATABASE'."

