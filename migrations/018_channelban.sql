create table channel_bans (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    user_id     bigint  not null,
    full_ban    boolean not null default false,

    primary key (channel_id, user_id)
);

create index channel_bans_server_member_idx on channel_bans (server_id, user_id);

---- create above / drop below ----

drop table channel_bans;