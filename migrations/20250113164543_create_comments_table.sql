-- +goose Up
CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          photo_id INT NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
                          user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);


-- +goose Down
DROP TABLE comments;

