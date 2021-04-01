create table pending_actions (
    scheduled   timestamp   not null,

    user_id bigint,
    server_id bigint,

    add_roles       bigint[]    not null    default array[]::bigint[],
    remove_roles    bigint[]    not null    default array[]::bigint[],

    message         json        not null    default '{}'
);

---- create above / drop below ----

drop table pending_actions;