create table mod_log (
    id          bigint,
    server_id   bigint,

    user_id     bigint  not null,

    action_type text    not null,
    reason      text    not null,

    time    timestamp   not null    default (current_timestamp at time zone 'utc'),

    unique(server_id, user_id)
);

create index mod_log_user_idx on mod_log (user_id, server_id);

create table server_mod_settings (
    id  bigint primary key,

    mute_role   bigint  not null    default 0
);

---- create above / drop below ----

drop table mod_log;
drop table server_mod_settings;