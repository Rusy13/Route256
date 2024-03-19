-- +goose Up
-- +goose StatementBegin
create table pvz(
      id BIGSERIAL PRIMARY KEY NOT NULL ,
      pvzname TEXT NOT NULL DEFAULT '',
      address TEXT NOT NULL DEFAULT '',
      email TEXT NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table pvz;
-- +goose StatementEnd

