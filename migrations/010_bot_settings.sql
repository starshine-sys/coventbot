create type user_status as enum ('online', 'idle', 'dnd');

create table bot_settings (
    id int primary key not null default 1, -- enforced only equal to 1

    status          user_status not null    default 'online',
    activity_type   text        not null    default 'playing',
    activity        text        not null    default '',

    constraint singleton check (id = 1)
);

insert into bot_settings (id, status) values (1, 'online');

---- create above / drop below ----

drop table bot_settings;
drop type user_status;