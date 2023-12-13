-- +migrate Up
create table if not exists proofs(
    id serial primary key,
    creator text not null,
    created_at timestamp without time zone not null default now(),
    proof jsonb not null
);

-- +migrate Down
drop table if exists proofs;
