-- +migrate Up
create table if not exists proofs(
    id         uuid primary key not null,
    creator    text not null,
    created_at timestamp without time zone not null default now(),
    proof      jsonb not null,
    type       text not null
);

create table if not exists links(
    id         uuid primary key not null,
    user_id    text             not null,
    created_at timestamp        not null default current_timestamp
);

create table if not exists links_to_proofs(
    link_id  uuid not null,
    proof_id uuid not null,
    primary key (link_id, proof_id),
    foreign key (link_id) references links(id) on delete cascade,
    foreign key (proof_id) references proofs(id) on delete cascade
);

-- +migrate Down
drop table if exists proofs;
drop table if exists links;
drop table if exists links_to_proofs;
