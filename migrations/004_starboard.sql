create table starboard_messages
(
    message_id              bigint primary key,
    channel_id              bigint not null,
    server_id               bigint not null,
    starboard_message_id    bigint
);

create index starboard_messages_id_idx on starboard_messages (starboard_message_id);

alter table servers add column starboard_channel bigint not null default 0;
alter table servers add column starboard_emoji text not null default '⭐';
alter table servers add column starboard_limit int not null default 3;

alter table servers add column starboard_blacklist bigint[] not null default array[]::bigint[];

---- create above / drop below ----

drop table starboard_messages;
alter table servers drop column starboard_channel;
alter table servers drop column starboard_emoji;
alter table servers drop column starboard_limit;
alter table servers drop column starboard_blacklist;