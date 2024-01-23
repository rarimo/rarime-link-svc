-- +migrate Up
alter table links_to_proofs drop constraint links_to_proofs_link_id_fkey;
alter table links alter column id set data type text;
alter table links_to_proofs alter column link_id set data type text;
alter table links_to_proofs add foreign key (link_id) references links(id) on delete cascade;
