begin;

create table load_test (
    id uuid primary key,
    value text not null default ''
);

commit;