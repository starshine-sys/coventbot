-- +migrate Up

-- 2022-03-13
-- Add webhooks for starboard

create table starboard_webhooks (
    id          bigint  primary key,
    channel_id  bigint  unique  not null,
    token       text
);

alter table starboard_messages add column webhook_id bigint references starboard_webhooks (id) on delete set null;

alter table servers add column starboard_avatar_url text not null default '';
alter table servers add column starboard_username text not null default '';
