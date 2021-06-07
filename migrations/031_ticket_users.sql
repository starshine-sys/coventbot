alter table tickets add column users bigint[] not null default array[]::bigint[];

update tickets set users = array[owner_id]::bigint[] where users = array[]::bigint[];

---- create above / drop below ----

alter table tickets drop column users;