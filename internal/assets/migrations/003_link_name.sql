-- +migrate Up
alter table links add column name text unique;

-- +migrate Down
alter table links drop column name;
