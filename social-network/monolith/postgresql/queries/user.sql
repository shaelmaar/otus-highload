-- name: UserCreate :exec
insert into "user"(id, password_hash, first_name, second_name,
                   birth_date, gender, biography, city)
values ( @id, @password_hash, @first_name, @second_name,
         @birth_date, @gender, @biography, @city);

-- name: UserGetByID :one
select * from "user"
where id = @id;

-- name: UsersGetByFirstNameSecondName :many
select * from "user"
where first_name ilike '%' || @first_name::text || '%'
    and second_name ilike '%' || @second_name::text || '%'
order by id;

-- name: UsersMassCreate :copyfrom
insert into "user"(id, password_hash, first_name, second_name,
                   birth_date, gender, biography, city)
values ($1, $2, $3, $4,
        $5, $6, $7, $8);
