-- +migrate Up

-- 2022-04-07
-- Clean up some now-unused tables

-- username/nickname logging
drop table usernames;
drop table nicknames;

-- channel mirror
drop table channel_mirror_messages;
drop table channel_mirror;

-- gatekeeper
drop table gatekeeper;
alter table servers drop column welcome_channel;
alter table servers drop column welcome_message;
alter table servers drop column member_role;

-- tags
drop table tags;
alter table servers drop column tag_mod_role;
