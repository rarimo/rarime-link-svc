-- +migrate Up
create table if not exists links(
  id serial primary key,
  index integer not null,
  created_at timestamp not null default current_timestamp
);

create table  if not exists links_to_proofs(
  id serial primary key,
  link_id integer not null,
  foreign key (link_id) references links(id) on delete cascade
);
-- +migrate Down
drop table if exists links;
drop table if exists links_to_proofs;
