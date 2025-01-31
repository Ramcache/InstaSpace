-- +goose Up
CREATE TABLE photo_likes (
                             id SERIAL PRIMARY KEY,
                             photo_id INT NOT NULL,
                             user_id INT NOT NULL,
                             created_at TIMESTAMP DEFAULT NOW(),
                             CONSTRAINT fk_photo FOREIGN KEY (photo_id) REFERENCES photos(id),
                             CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
                             UNIQUE (photo_id, user_id)
);

-- +goose Down
DROP TABLE photo_likes;
