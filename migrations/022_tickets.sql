create table ticket_categories (
    category_id bigint  primary key,

    server_id   bigint  not null,

    per_user_limit  integer     not null    default -1,
    log_channel     bigint      not null,
    count           integer     not null    default 0,

    can_creator_close   boolean not null    default true,

    mention     text    not null    default '',
    description text    not null    default ''
);

create table tickets (
    channel_id  bigint  primary key,

    category_id bigint  not null references ticket_categories (category_id) on delete cascade,

    owner_id    bigint  not null
);

---- create above / drop below ----

drop table tickets;
drop table ticket_categories;