-- name: LoadTestInsert :exec
insert into load_test(id, value)
values ( @id, @value);

-- name: LoadTestDelete :exec
delete from load_test
where id = @id;
