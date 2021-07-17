create table highlights (
    user_id     bigint,
    server_id   bigint,
    highlights  text[]      not null    default array[]::text[],
    blocked     bigint[]    not null    default array[]::bigint[],

    primary key (user_id, server_id)
);

create table highlight_config (
    server_id   bigint      primary key,
    hl_enabled  bool        not null    default false,
    blocked     bigint[]    not null    default array[]::bigint[]
);

create table highlight_delete_queue (
    message_id  bigint  primary key,
    channel_id  bigint  not null
);

---- create above / drop below ----

drop table highlights;
drop table highlight_config;
drop table highlight_delete_queue;