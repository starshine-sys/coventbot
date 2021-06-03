alter table servers add column keyrole_channel bigint not null default 0;
alter table servers add column keyroles bigint[] not null default array[]::bigint[];

---- create above / drop below ----

alter table servers drop column keyrole_channel;
alter table servers drop column keyroles;