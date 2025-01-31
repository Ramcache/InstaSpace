-- +goose Up
CREATE TABLE photos (
                        id SERIAL PRIMARY KEY,
                        user_id INT NOT NULL,
                        url TEXT NOT NULL,
                        description TEXT,
                        likes_count INT DEFAULT 0,
                        created_at TIMESTAMP DEFAULT NOW(),
                        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS photos;
