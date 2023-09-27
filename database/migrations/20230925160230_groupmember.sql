-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_members(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    uid        bigint NOT NULL,
    group_id   bigint NOT NULL,
    inviter    bigint NOT NULL,
    remark     char(64) NOT NULL DEFAULT '',
    status     smallint NOT NULL DEFAULT 0,
    apply_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table group_members;
-- +goose StatementEnd
