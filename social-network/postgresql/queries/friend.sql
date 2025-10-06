-- name: FriendCrete :exec
insert into friend(user_id, friend_id, created_at)
values (@user_id, @friend_id, @created_at)
on conflict (user_id, friend_id) do nothing;

-- name: FriendDelete :exec
delete
from friend
where user_id = @user_id and friend_id = @friend_id;

-- name: FriendIDsByUserID :many
select friend_id
from friend
where user_id = @user_id;

-- name: UserIDsByFriendID :many
select user_id
from friend
where friend_id = @friend_id;
