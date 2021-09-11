-- +migrate Up

-- 2021-09-11
-- Add timezone to user config, used for reminders

alter table user_config add column timezone text not null default '';