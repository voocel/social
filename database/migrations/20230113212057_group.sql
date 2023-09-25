-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups(
    id           BIGSERIAL PRIMARY KEY NOT NULL,
    name         varchar(32) NOT NULL DEFAULT '未命名',
    owner        bigint NOT NULL DEFAULT 0,
    created_uid  bigint NOT NULL DEFAULT 0,
    mode         smallint  NOT NULL DEFAULT 0,
    type         smallint  NOT NULL DEFAULT 0,
    status       smallint  NOT NULL DEFAULT 0,
    invite_mode  smallint  NOT NULL DEFAULT 0,
    notice       varchar(255) NOT NULL DEFAULT '',
    introduction varchar(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at   TIMESTAMP
);

COMMENT ON COLUMN groups.status IS '群组状态: 0 正常, 1 已解散';
COMMENT ON COLUMN groups.mode IS '禁言模式: 0 无, 1 仅管理员可发言, 2 全员禁言';
COMMENT ON COLUMN groups.invite_mode IS '邀请模式: 0 仅群主可邀请, 1 仅管理员可邀请, 2 所有成员可邀请';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table group;
-- +goose StatementEnd
