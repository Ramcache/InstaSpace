-- +goose Up
CREATE TABLE conversations (
                               id SERIAL PRIMARY KEY,
                               user1_id INT NOT NULL,
                               user2_id INT NOT NULL,
                               created_at TIMESTAMP DEFAULT NOW(),
                               FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
                               UNIQUE(user1_id, user2_id)
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          conversation_id INT NOT NULL,
                          sender_id INT NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW(),
                          FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
                          FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE conversations;

DROP TABLE messages;
