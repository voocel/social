-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       varchar(32) NOT NULL DEFAULT '未命名',
    owner      bigint NOT NULL DEFAULT 0,
    notice     varchar(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table group;
-- +goose StatementEnd
