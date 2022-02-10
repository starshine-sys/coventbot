-- +migrate Up

-- 2022-02-10
-- Add Lua-based custom commands

create table custom_commands (
    id          bigserial primary key,
    guild_id    bigint  not null,
    name        text    not null,
    source      text    not null
);

create unique index custom_commands_name_idx on custom_commands (guild_id, lower(name));
