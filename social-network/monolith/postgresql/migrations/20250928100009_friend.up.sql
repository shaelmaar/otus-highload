begin;

create table if not exists friend
(
    user_id    uuid                     not null references "user" (id) on delete cascade,
    friend_id  uuid                     not null references "user" (id) on delete cascade,
    created_at timestamp with time zone not null default now(),
    primary key (user_id, friend_id)
);

comment on table friend is 'информация о подписках (друзьях) пользователей';
comment on column friend.user_id is 'идентификатор пользователя подписчика';
comment on column friend.friend_id is 'идентификатор пользователя, на которого подписан пользователь';

create index if not exists friend__friend_id_idx on friend (friend_id);

commit;