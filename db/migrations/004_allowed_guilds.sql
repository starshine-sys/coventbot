-- +migrate Up

-- 2021-12-05
-- Add allowed guilds table, bot will leave any guild not in this table

create table allowed_guilds (
    id          bigint      primary key,
    reason      text        not null,
    added_by    bigint      not null,
    added_for   bigint      not null,
    added_at    timestamp   not null default (current_timestamp at time zone 'utc')
);
