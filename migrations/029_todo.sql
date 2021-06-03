create table todos (
    id  serial  primary key,

    user_id     bigint  not null,
    description text    not null,

    orig_mid        bigint  not null,
    orig_channel_id bigint  not null,
    orig_server_id  bigint  not null,

    mid         bigint  not null,
    channel_id  bigint  not null,
    server_id   bigint  not null,

    complete bool not null default false,

    created     timestamp   not null default (current_timestamp at time zone 'utc'),
    completed   timestamp   default null
);

create index todos_user_idx on todos (user_id);

alter table user_config add column todo_channel bigint not null default 0;

---- create above / drop below ----

drop table todos;
alter table server_levels drop column todo_channel;