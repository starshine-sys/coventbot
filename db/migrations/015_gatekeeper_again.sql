-- +migrate Up

-- 2022-08-22
-- Readd the gatekeeper table

create table gatekeeper (
    guild_id    bigint  not null,
    user_id     bigint  not null,
    key         uuid    not null,

    pending     boolean not null,

    unique(guild_id, user_id)
);

alter table servers add column welcome_channel not null default 0;
alter table servers add column welcome_message text not null default 'Welcome to the server, {mention}!';
alter table servers add column member_role not null default 0;
