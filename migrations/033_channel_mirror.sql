create table channel_mirror (
    server_id       bigint not null,

    from_channel    bigint  primary key,
    to_channel      bigint  not null,

    webhook_id  bigint  not null,
    token       text    not null
);

create table channel_mirror_messages (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  primary key,
    original    bigint  not null,
    user_id     bigint  not null
);

---- create above / drop below ----

drop table channel_mirror;
drop table channel_mirror_messages;