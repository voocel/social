-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friend_applies(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    from_id    bigint NOT NULL,
    to_id      bigint NOT NULL,
    remark     char(64) NOT NULL DEFAULT '',
    status     smallint NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP
);

COMMENT ON COLUMN friend_applies.from_id IS '申请人ID';
COMMENT ON COLUMN friend_applies.to_id IS '被申请人ID';
COMMENT ON COLUMN friend_applies.remark IS '备注';
COMMENT ON COLUMN friend_applies.status IS '申请状态: 0 发起申请, 1 已通过, 2 已拒绝';
COMMENT ON COLUMN friend_applies.created_at IS '创建时间';
COMMENT ON COLUMN friend_applies.updated_at IS '修改时间';
COMMENT ON COLUMN friend_applies.deleted_at IS '删除时间';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table friend_apply;
-- +goose StatementEnd
