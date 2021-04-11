create table user_config (
    user_id bigint primary key,

    use_enko_format boolean not null    default false,

    opt_out_nicknames   boolean not null    default false,
    opt_out_usernames   boolean not null    default false
);

---- create above / drop below ----

drop table user_config;