create type level_messages as enum ('ALL_DM', 'REWARDS_DM', 'ALL_CHANNEL', 'REWARDS_CHANNEL', 'NONE');

alter table server_levels add column level_messages level_messages not null default 'NONE';
alter table server_levels add column level_channel bigint not null default 0;

---- create above / drop below ----

alter table server_levels drop column level_channel;
alter table server_levels drop column level_messages;
drop type level_messages;