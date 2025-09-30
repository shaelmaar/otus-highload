-- name: PostCreate :exec
insert into post(id, content, author_user_id, created_at)
values (@id, @content, @author_user_id, @created_at);

-- name: PostGetByID :one
select *
from post
where id = @id;

-- name: PostUpdate :exec
update post
set content = @content,
    updated_at = @updated_at
where id = @id;

-- name: PostGetWithLockByID :one
select *
from post
where id = @id
for no key update;

-- name: PostDelete :exec
delete
from post
where id = @id;

-- name: LastPostsByUserIDsWithOffsetLimit :many
select *
from post
where author_user_id = any((@user_ids::text[])::uuid[])
order by created_at desc
limit sqlc.arg('limit')
offset sqlc.arg('offset');
