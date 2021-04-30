create table roles (
    id          serial  primary key,
    server_id   bigint  not null,
    name        text    not null,

    require_role    bigint      not null    default 0,
    roles           bigint[]    not null    default array[]::bigint[]
);

create index roles_server_id_idx on roles (server_id);

---- create above / drop below ----

drop table roles;