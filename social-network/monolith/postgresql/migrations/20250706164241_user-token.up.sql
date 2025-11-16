begin;

create table if not exists user_token (
    id              bigserial primary key,
    user_id         uuid not null references "user"(id) on delete cascade,
    token           text not null unique,
    expires_at      timestamp with time zone,
    created_at      timestamp with time zone not null default now()
);

comment on table user_token is 'токен пользователя';
comment on column user_token.user_id is 'идентификатор пользователя';
comment on column user_token.token is 'токен';
comment on column user_token.expires_at is 'время срока годности токена';

commit;