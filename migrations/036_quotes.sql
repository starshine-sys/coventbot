create table quotes (
    id  serial  primary key,
    hid char(5) not null,

    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  not null,
    user_id     bigint  not null,
    added_by    bigint  not null,
    content     text    not null,

    added   timestamp   not null default (current_timestamp at time zone 'utc'),

    unique (server_id, hid)
);

create function generate_hid() returns char(5) as $$
    select string_agg(substr('abcdefghijklmnopqrstuvwxyz', ceil(random() * 26)::integer, 1), '') from generate_series(1, 5)
$$ language sql volatile;

create function find_free_quote_hid(guild bigint) returns char(5) as $$
declare new_hid char(5);
begin
    loop
        new_hid := generate_hid();
        if not exists (select 1 from quotes where hid = new_hid and server_id = guild) then return new_hid; end if;
    end loop;
end
$$ language plpgsql volatile;

alter table servers add column quotes_enabled boolean not null default false;
alter table servers add column quote_suppress_messages boolean not null default false;

---- create above / drop below ----

drop function find_free_quote_hid;
drop function generate_hid;
drop table quotes;

alter table servers drop column quotes_enabled;
alter table servers drop column quote_suppress_messages;