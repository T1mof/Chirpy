-- +goose Up
CREATE TABLE refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id uuid NOT NULL,
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;