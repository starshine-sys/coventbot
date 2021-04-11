create table pk_messages (
    msg_id      bigint  primary key,
    user_id     bigint  not null    default 0,
    channel_id  bigint  not null    default 0,
    server_id   bigint  not null    default 0,

    username    text    not null,
    member      text    not null,
    system      text    not null,

    content     text    not null
);

alter table servers add column pk_log_channel bigint not null default 0;

---- create above / drop below ----

drop table pk_messages;
alter table servers drop column pk_log_channel;