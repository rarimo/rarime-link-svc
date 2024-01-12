-- +migrate Up
alter table proofs add column schema_url text not null default '';

-- +migrate Down
alter table proofs drop column schema_url;
