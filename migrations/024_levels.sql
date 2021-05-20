create table server_levels (
    id  bigint  primary key,

    blocked_channels    bigint[]    not null    default array[]::bigint[],
    blocked_roles       bigint[]    not null    default array[]::bigint[],
    blocked_categories  bigint[]    not null    default array[]::bigint[],

    between_xp  interval    not null    default '1 minute',
    reward_text text        not null    default '',

    levels_enabled          boolean not null    default true,
    leaderboard_mod_only    boolean not null    default false,
    show_next_reward        boolean not null    default true
);

create table levels (
    server_id   bigint  not null,
    user_id     bigint  not null,

    xp      bigint  not null    default 0,

    next_time   timestamp   not null    default (current_timestamp at time zone 'utc'),

    primary key (server_id, user_id)
);

create table level_rewards (
    server_id   bigint  not null,
    lvl         bigint  not null,
    role_reward bigint  not null,

    primary key (server_id, lvl)
);

---- create above / drop below ----

drop table server_levels;
drop table levels;
drop table level_rewards;