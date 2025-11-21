begin;

create table if not exists post
(
    id             uuid                     not null primary key,
    content        text                     not null default '',
    author_user_id uuid                     not null references "user" (id) on delete cascade,
    created_at     timestamp with time zone not null default now(),
    updated_at timestamp with time zone
);

comment on table post is 'посты пользователей';
comment on column post.id is 'идентификатор поста';
comment on column post.content is 'содержимое поста';
comment on column post.author_user_id is 'идентификатор пользователя автора поста';

create index if not exists post__author_user_id_created_at_desc_idx on post (author_user_id, created_at desc);

commit;