-- +migrate Up
alter table proofs add column field text not null;

-- +migrate Down
alter table proofs drop column field;
