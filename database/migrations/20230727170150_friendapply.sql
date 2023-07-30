-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friend_apply(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    from_id    bigint NOT NULL,
    to_id      bigint NOT NULL,
    remark     char(64) NOT NULL DEFAULT '',
    status     smallint NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friend_apply;
-- +goose StatementEnd
