create table tags (
    id          int     not null default 0,
    server_id   bigint  not null default 0,

    name        text    not null default '',
    response    text    not null default '',

    created_by  bigint      not null default 0,
    created_at  timestamp   not null default (current_timestamp at time zone 'utc')
);

create unique index tags_uniq_idx on tags (lower(name), server_id);

create table servers (
    id  bigint  primary key,

    prefixes text[] not null default array[]::text[],

    blacklisted_channels bigint[] not null default array[]::bigint[]
);

---- create above / drop below ----

drop table tags;
drop table servers;
