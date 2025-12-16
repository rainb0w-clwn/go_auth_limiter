-- +goose Up
-- +goose StatementBegin
insert into rate_limit(type, value, description)
values ('login', 10, 'Ограничение для логина'),
       ('password', 100, 'Ограничение для пароля'),
       ('ip', 1000, null);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete
from rate_limit
where type in ('login', 'password', 'ip');
-- +goose StatementEnd
