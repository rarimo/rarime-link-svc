-- +migrate Up
create type proof_operator as enum('$noop', '$eq', '$lt', '$gt', '$in', '$nin', '$ne');
alter table proofs add column operator proof_operator not null;

-- +migrate Down
alter table proofs drop column operator;
drop type proof_operator;
