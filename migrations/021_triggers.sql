create table triggers (
    message_id bigint,
    emoji text,

    command text[] not null default array[]::text[],

    primary key (message_id, emoji)
);

---- create above / drop below ----

drop table triggers;