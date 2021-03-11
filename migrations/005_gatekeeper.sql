create table gatekeeper (
    server_id   bigint  not null,
    user_id     bigint  not null,
    key         uuid    not null,

    pending     boolean not null,

    unique(server_id, user_id)
);

alter table servers add column member_role bigint not null default 0;
alter table servers add column welcome_channel bigint not null default 0;
alter table servers add column welcome_message text not null default 'Welcome to {server}, {mention}!';

---- create above / drop below ----

drop table gatekeeper;
alter table servers drop column member_role;
alter table servers drop column welcome_channel;
alter table servers drop column welcome_message;