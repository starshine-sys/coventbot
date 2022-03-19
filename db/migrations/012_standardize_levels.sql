-- +migrate Up

-- 2022-01-24
-- Drop carline_curve config setting, all servers now use Carl-bot's level curve.

alter table server_levels drop column carline_curve;
