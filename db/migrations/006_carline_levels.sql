-- +migrate Up

-- 2022-01-24
-- Add voice XP

alter table server_levels add column carline_curve boolean not null default false;
