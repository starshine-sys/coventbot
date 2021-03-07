alter table servers add column watch_list_channel bigint not null default 0;
alter table servers add column watch_list bigint[] not null default array[]::bigint[];

---- create above / drop below ----

alter table servers drop column watch_list_channel;
alter table servers drop column watch_list;