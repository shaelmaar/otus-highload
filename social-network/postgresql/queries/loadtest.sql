-- name: LoadTestInsert :exec
insert into load_test(id, value)
values ( @id, @value);