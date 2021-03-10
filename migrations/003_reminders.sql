create table reminders (
    id          serial  primary key,
    user_id     bigint  not null default 0,
    message_id  bigint  not null default 0,
    channel_id  bigint  not null default 0,
    server_id   bigint  not null default 0,

    reminder    text    not null default '',

    expires timestamp not null
);

create index reminders_user_idx on reminders (user_id);
create index reminders_user_server_idx on reminders (user_id, server_id);

---- create above / drop below ----

drop table reminders;