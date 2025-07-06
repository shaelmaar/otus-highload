begin;

create type gender as enum ('male', 'female');
comment on type gender is 'Пол:
* male          - мужской.
* female        - женский.
';

create table if not exists "user" (
    id uuid not null primary key,
    password_hash text not null,
    first_name text not null,
    second_name text not null,
    birth_date date not null,
    gender gender not null,
    biography text not null default '',
    city text not null default '',
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone
);

comment on table "user" is 'Пользователь';
comment on column  "user".password_hash is 'Хэш пароля';
comment on column  "user".first_name is 'Имя';
comment on column  "user".second_name is 'Фамилия';
comment on column  "user".birth_date is 'Дата рождения';
comment on column  "user".gender is 'Пол';
comment on column  "user".biography is 'Хобби, интересы и т.п.';
comment on column  "user".city is 'Город';


commit;
