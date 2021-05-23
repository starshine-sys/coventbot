create table react_roles (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  primary key
);

create table react_role_entries (
    message_id  bigint  not null    references react_roles (message_id) on delete cascade,
    emote       text    not null, -- this can either be a default emoji or a custom emote ID
    role_id     bigint  not null,

    primary key (message_id, emote)
);

create index react_roles_server_idx on react_roles (server_id);
create index react_role_entries_message_id_idx on react_role_entries (message_id);

---- create above / drop below ----

drop table react_roles;
drop table react_role_entries;