-- +migrate Up

-- 2022-01-17
-- Add voice XP

alter table server_levels add column voice boolean not null default false;
