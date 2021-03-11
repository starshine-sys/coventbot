alter table servers add column approve_remove_roles bigint[] not null default array[]::bigint[];
alter table servers add column approve_add_roles bigint[] not null default array[]::bigint[];
alter table servers add column approve_welcome_channel bigint not null default 0;
alter table servers add column approve_welcome_message text not null default 'Welcome to {server}, {mention}!';

---- create above / drop below ----

alter table servers drop column approve_remove_roles;
alter table servers drop column approve_add_roles;
alter table servers drop column approve_welcome_channel;
alter table servers drop column approve_welcome_message;