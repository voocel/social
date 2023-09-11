-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friend(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    uid        bigint NOT NULL,
    friend_id  bigint NOT NULL,
    remark     char(64) NOT NULL DEFAULT '',
    shield     smallint NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP
)

-- CREATE UNIQUE INDEX idx_uid_friend_id ON friend (uid, friend_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friend;
-- +goose StatementEnd
