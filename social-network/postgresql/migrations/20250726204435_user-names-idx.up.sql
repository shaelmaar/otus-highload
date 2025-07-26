begin;

create index if not exists user_names_gin_trgm_idx on "user"
    using gin (first_name gin_trgm_ops, second_name gin_trgm_ops);

commit;