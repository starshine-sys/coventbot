alter table server_levels add column reward_log bigint not null default 0;
alter table server_levels add column nolevels_log bigint not null default 0;

create table nolevels (
    server_id   bigint  not null,
    user_id     bigint  not null,

    expires boolean     not null    default false,
    -- this default is not used if "expiry" is also left as the default so it's fine
    expiry  timestamp   not null    default (current_timestamp at time zone 'utc'),

    primary key (server_id, user_id)
);

create index nolevels_server_idx on nolevels (server_id);

---- create above / drop below ----

alter table server_levels drop column reward_log;
alter table server_levels drop column nolevels_log;
drop table nolevels;