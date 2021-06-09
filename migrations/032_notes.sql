create table notes (
    id  serial  primary key,

    server_id   bigint not null,
    user_id     bigint not null,

    note        text        not null,
    moderator   bigint      not null,
    created     timestamp   not null default (current_timestamp at time zone 'utc')
);

create index notes_server_idx on todos (server_id);
create index notes_server_user_idx on todos (server_id, user_id);

---- create above / drop below ----

drop table notes;