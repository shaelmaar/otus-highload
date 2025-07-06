-- name: UserCreate :exec
insert into "user"(id, password_hash, first_name, second_name,
                   birth_date, gender, biography, city)
values ( @id, @password_hash, @first_name, @second_name,
         @birth_date, @gender, @biography, @city);

-- name: UserGetByID :one
select * from "user"
where id = @id;

-- name: UserTokenCreate :one
insert into user_token(user_id, token, expires_at)
values (@user_id, @token, @expires_at)
returning id;

-- name: UserTokenDeleteByUserID :exec
delete from user_token
where user_id = @user_id;
