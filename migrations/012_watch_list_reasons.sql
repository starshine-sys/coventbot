create table watch_list_reasons (
    user_id     bigint,
    server_id   bigint,

    reason text not null default '',

    primary key (user_id, server_id)
);

---- create above / drop below ----

drop table watch_list_reasons;