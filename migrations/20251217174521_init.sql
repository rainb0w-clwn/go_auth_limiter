-- +goose Up
-- +goose StatementBegin
create table ip_net_rule
(
    id   bigint generated always as identity primary key,
    ip   varchar(50)  not null,
    type varchar(255) not null
);
-- +goose StatementEnd
-- +goose StatementBegin
create table rate_limit
(
    type        varchar(50) primary key,
    value       int          not null,
    description varchar(255) null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ip_net_rule;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS rate_limit;
-- +goose StatementEnd
