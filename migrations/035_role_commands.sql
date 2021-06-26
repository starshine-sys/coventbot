create table role_categories (
    id          serial  primary key,
    server_id   bigint  not null,

    name        text    not null default '',
    description text    not null default '',
    colour      int     not null default 0,

    require_role    bigint      not null    default 0,
    roles           bigint[]    not null    default array[]::bigint[],

    unique (server_id, name)
);

create index role_categories_server_idx on role_categories (server_id);

---- create above / drop below ----

drop table role_categories;