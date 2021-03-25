-- usernames contains the username log, username + discriminator
create table usernames (
    user_id bigint      not null,
    time    timestamp   not null default (current_timestamp at time zone 'utc'),
    name    text        not null
);

-- nicknames contains the nickname log
create table nicknames (
    server_id   bigint  not null,
    user_id     bigint  not null,
    time    timestamp   not null default (current_timestamp at time zone 'utc'),
    name    text        not null
);

create index usernames_idx on usernames (user_id);
create index nicknames_idx on nicknames (server_id, user_id);

---- create above / drop below ----

drop table usernames;
drop table nicknames;
