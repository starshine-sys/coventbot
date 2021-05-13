create table slowmode (
    server_id   bigint not null default 0,
    channel_id  bigint   primary key,
    slowmode    interval not null
);

create table user_slowmode (
    server_id   bigint      not null,
    channel_id  bigint      not null references slowmode (channel_id) on delete cascade,
    user_id     bigint      not null,
    expiry      timestamp   not null,

    primary key (channel_id, user_id)
);

alter table servers add column slowmode_ignore_role bigint not null default 0;

---- create above / drop below ----

drop table slowmode;
drop table user_slowmode;
alter table servers drop column slowmode_ignore_role;