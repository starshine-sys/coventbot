alter table servers add column mod_log_channel bigint not null default 0;

drop table mod_log;

create table mod_log (
    id          bigint,
    server_id   bigint,

    user_id     bigint  not null,
    mod_id      bigint  not null,

    action_type text    not null,
    reason      text    not null,

    time    timestamp   not null    default (current_timestamp at time zone 'utc'),

    unique(id, server_id)
);