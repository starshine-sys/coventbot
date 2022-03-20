-- +migrate Up

-- 2022-01-24
-- Add Carl-bot compatible level curve

alter table server_levels add column carline_curve boolean not null default false;
