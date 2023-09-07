-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id              BIGSERIAL PRIMARY KEY NOT NULL,
    username        varchar(32) NOT NULL DEFAULT '未知',
    password        char(64) NOT NULL DEFAULT '',
    mobile          char(11) NOT NULL DEFAULT '',
    nickname        varchar(50) NOT NULL DEFAULT '',
    email           varchar(50) NOT NULL DEFAULT '',
    avatar          varchar(128) NOT NULL DEFAULT '',
    summary         varchar(128) NOT NULL DEFAULT '',
    sex             smallint  NOT NULL DEFAULT 0,
    status          smallint  NOT NULL DEFAULT 1,
    birthday        DATE,
    last_login_time TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at      TIMESTAMP
);

-- CREATE INDEX idx_account_mobile_email ON users (username, mobile, email);

COMMENT ON COLUMN users.username IS '用户名称';
COMMENT ON COLUMN users.password IS '用户密码';
COMMENT ON COLUMN users.mobile IS '用户手机号';
COMMENT ON COLUMN users.nickname IS '用户昵称';
COMMENT ON COLUMN users.email IS '电子邮箱';
COMMENT ON COLUMN users.avatar IS '用户头像';
COMMENT ON COLUMN users.summary IS '用户简介';
COMMENT ON COLUMN users.sex IS '0未知，1男，2女';
COMMENT ON COLUMN users.status IS '用户账户有效状态，0正常1无效';
COMMENT ON COLUMN users.birthday IS '用户出生日期，一般年月日即可';
COMMENT ON COLUMN users.last_login_time IS '最后登录时间';
COMMENT ON COLUMN users.created_at IS '用户创建时间';
COMMENT ON COLUMN users.updated_at IS '用户修改时间';
COMMENT ON COLUMN users.deleted_at IS '用户删除时间';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
