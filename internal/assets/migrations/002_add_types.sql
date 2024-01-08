-- +migrate Up
alter table proofs
    add column type text;

-- +migrate Down
alter table proofs
    drop column type;
