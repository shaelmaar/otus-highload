-- name: FriendCrete :exec
insert into friend(user_id, friend_id, created_at)
values (@user_id, @friend_id, @created_at)
on conflict (user_id, friend_id) do nothing;

-- name: FriendDelete :exec
delete
from friend
where user_id = @user_id and friend_id = @friend_id;
