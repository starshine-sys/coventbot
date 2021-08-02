-- +migrate Up

-- 2021-08-02
-- Add channel and message IDs to the mod log for editing mod log reasons.
alter table mod_log add column channel_id bigint not null default 0;
alter table mod_log add column message_id bigint not null default 0;