-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages(
    id           BIGSERIAL PRIMARY KEY NOT NULL,
    sender_id    bigint NOT NULL,
    receiver_id  bigint NOT NULL,
    content      char(256) NOT NULL DEFAULT '',
    content_type smallint NOT NULL DEFAULT 0,
    status       smallint NOT NULL DEFAULT 0,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
