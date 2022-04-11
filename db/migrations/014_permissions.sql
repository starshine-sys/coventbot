-- +migrate Up

-- 2022-04-11
-- Overhaul the permissions system
-- Add a permission nodes table and rename the current permission levels
-- helper ⇒ moderator, moderator ⇒ manager

create table permission_nodes (
    guild_id    bigint  not null    references servers (id) on delete cascade,
    name        text    not null,
    level       int     not null,

    primary key (guild_id, name)
);

alter table servers rename column helper_roles to moderator_roles;
alter table servers rename column mod_roles to manager_roles;
