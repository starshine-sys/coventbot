-- started this a while back and then abandoned it, so gotta clean up first
drop table user_config;

create table user_config (
    user_id     bigint  primary key,

    disable_levelup_messages    boolean not null    default false,
    reminders_in_dm             boolean not null    default false,
    usernames_opt_out           boolean not null    default false
);

---- create above / drop below ----

drop table user_config;