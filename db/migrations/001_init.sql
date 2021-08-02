-- +migrate Up

-- 2021-08-02
-- Initial database schema, makes sure a database is up-to-date to the last tern migration.
-- Should be idempotent if the database is already updated to the latest migration.

-- +migrate StatementBegin

-- bot status + activity
do $$ begin
    create type user_status as enum ('online', 'idle', 'dnd');
exception when duplicate_object then null;
end $$;
-- +migrate StatementEnd

-- +migrate StatementBegin

-- level message setting
do $$ begin
    create type level_messages as enum ('ALL_DM', 'REWARDS_DM', 'ALL_CHANNEL', 'REWARDS_CHANNEL', 'NONE');
exception when duplicate_object then null;
end $$;
-- +migrate StatementEnd

create table if not exists tags (
    id          int     not null default 0,
    server_id   bigint  not null default 0,

    name        text    not null default '',
    response    text    not null default '',

    created_by  bigint      not null default 0,
    created_at  timestamp   not null default (current_timestamp at time zone 'utc')
);

-- server configuration
create table if not exists servers (
    id  bigint  primary key,

    prefixes text[] not null default array[]::text[],

    helper_roles bigint[] not null default array[]::bigint[],
    mod_roles bigint[] not null default array[]::bigint[],
    admin_roles bigint[] not null default array[]::bigint[],
    roles_set_up boolean not null default false,

    blacklisted_channels bigint[] not null default array[]::bigint[],
    watch_list_channel bigint not null default 0,
    watch_list bigint[] not null default array[]::bigint[],
    starboard_channel bigint not null default 0,
    starboard_emoji text not null default '⭐',
    starboard_limit int not null default 3,
    starboard_blacklist bigint[] not null default array[]::bigint[],
    member_role bigint not null default 0,
    welcome_channel bigint not null default 0,
    welcome_message text not null default 'Welcome to {server}, {mention}!',
    approve_remove_roles bigint[] not null default array[]::bigint[],
    approve_add_roles bigint[] not null default array[]::bigint[],
    approve_welcome_channel bigint not null default 0,
    approve_welcome_message text not null default 'Welcome to {server}, {mention}!',
    tag_mod_role bigint not null default 0,
    pk_log_channel bigint not null default 0,
    slowmode_ignore_role bigint not null default 0,
    mod_log_channel bigint not null default 0,
    keyrole_channel bigint not null default 0,
    keyroles bigint[] not null default array[]::bigint[],
    quotes_enabled boolean not null default false,
    quote_suppress_messages boolean not null default false,

    mute_role bigint not null default 0,
    pause_role bigint not null default 0,
    muteme_message text not null default 'Successfully {action}d {mention} for {duration}.'
);

-- starboard messages
create table if not exists starboard_messages (
    message_id              bigint primary key,
    channel_id              bigint not null,
    server_id               bigint not null,
    starboard_message_id    bigint
);

-- gatekeeper keys
create table if not exists gatekeeper (
    server_id   bigint  not null,
    user_id     bigint  not null,
    key         uuid    not null,

    pending     boolean not null,

    unique(server_id, user_id)
);

-- usernames contains the username log, username + discriminator
create table if not exists usernames (
    user_id bigint      not null,
    time    timestamp   not null default (current_timestamp at time zone 'utc'),
    name    text        not null
);

-- nicknames contains the nickname log
create table if not exists nicknames (
    server_id   bigint  not null,
    user_id     bigint  not null,
    time    timestamp   not null default (current_timestamp at time zone 'utc'),
    name    text        not null
);

create table if not exists bot_settings (
    id int primary key not null default 1, -- enforced only equal to 1

    status          user_status not null    default 'online',
    activity_type   text        not null    default 'playing',
    activity        text        not null    default '',

    constraint singleton check (id = 1)
);

create table if not exists server_mod_settings (
    id  bigint primary key,

    mute_role   bigint  not null    default 0
);

create table if not exists watch_list_reasons (
    user_id     bigint,
    server_id   bigint,

    reason text not null default '',

    primary key (user_id, server_id)
);

create table if not exists pk_messages (
    msg_id      bigint  primary key,
    user_id     bigint  not null    default 0,
    channel_id  bigint  not null    default 0,
    server_id   bigint  not null    default 0,

    username    text    not null,
    member      text    not null,
    system      text    not null,

    content     text    not null
);

create table if not exists roles (
    id          serial  primary key,
    server_id   bigint  not null,
    name        text    not null,

    require_role    bigint      not null    default 0,
    roles           bigint[]    not null    default array[]::bigint[]
);

create table if not exists slowmode (
    server_id   bigint not null default 0,
    channel_id  bigint   primary key,
    slowmode    interval not null
);

create table if not exists user_slowmode (
    server_id   bigint      not null,
    channel_id  bigint      not null references slowmode (channel_id) on delete cascade,
    user_id     bigint      not null,
    expiry      timestamp   not null,

    primary key (channel_id, user_id)
);

create table if not exists channel_bans (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    user_id     bigint  not null,
    full_ban    boolean not null default false,

    primary key (channel_id, user_id)
);

create table if not exists mod_log (
    id          bigint,
    server_id   bigint,

    user_id     bigint  not null,
    mod_id      bigint  not null,

    action_type text    not null,
    reason      text    not null,

    time    timestamp   not null    default (current_timestamp at time zone 'utc'),

    unique(id, server_id)
);

create table if not exists reminders (
    id          serial  primary key,
    user_id     bigint  not null default 0,
    message_id  bigint  not null default 0,
    channel_id  bigint  not null default 0,
    server_id   bigint  not null default 0,

    reminder    text    not null default '',

    set_time    timestamp   not null default (current_timestamp at time zone 'utc'),
    expires     timestamp not null
);

create table if not exists triggers (
    message_id bigint,
    emoji text,

    command text[] not null default array[]::text[],

    primary key (message_id, emoji)
);

create table if not exists ticket_categories (
    category_id bigint  primary key,

    server_id   bigint  not null,

    per_user_limit  integer     not null    default -1,
    log_channel     bigint      not null,
    count           integer     not null    default 0,

    can_creator_close   boolean not null    default true,

    name        text    not null    default '',
    mention     text    not null    default '',
    description text    not null    default ''
);

create table if not exists tickets (
    channel_id  bigint  primary key,

    category_id bigint  not null references ticket_categories (category_id) on delete cascade,

    owner_id    bigint      not null,
    users       bigint[]    not null    default array[]::bigint[]
);

create table if not exists server_levels (
    id  bigint  primary key,

    blocked_channels    bigint[]    not null    default array[]::bigint[],
    blocked_roles       bigint[]    not null    default array[]::bigint[],
    blocked_categories  bigint[]    not null    default array[]::bigint[],

    between_xp  interval    not null    default '1 minute',
    reward_text text        not null    default '',

    levels_enabled          boolean not null    default true,
    leaderboard_mod_only    boolean not null    default false,
    show_next_reward        boolean not null    default true,

    reward_log      bigint  not null    default 0,
    nolevels_log    bigint  not null    default 0,

    background text not null default '',

    level_messages level_messages not null default 'NONE',
    level_channel bigint not null default 0
);

create table if not exists levels (
    server_id   bigint  not null,
    user_id     bigint  not null,

    xp          bigint  not null    default 0,
    colour      bigint  not null    default 0,
    background  text    not null    default '',

    next_time   timestamp   not null    default (current_timestamp at time zone 'utc'),

    primary key (server_id, user_id)
);

create table if not exists level_rewards (
    server_id   bigint  not null,
    lvl         bigint  not null,
    role_reward bigint  not null,

    primary key (server_id, lvl)
);

create table if not exists nolevels (
    server_id   bigint  not null,
    user_id     bigint  not null,

    expires boolean     not null    default false,
    -- this default is not used if "expires" is also left as the default so it's fine
    expiry  timestamp   not null    default (current_timestamp at time zone 'utc'),

    primary key (server_id, user_id)
);

create table if not exists react_roles (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  primary key,

    title       text,
    description text,
    mention     boolean
);

create table if not exists react_role_entries (
    message_id  bigint  not null    references react_roles (message_id) on delete cascade,
    emote       text    not null, -- this can either be a default emoji or a custom emote ID
    role_id     bigint  not null,

    primary key (message_id, emote)
);

create table if not exists user_config (
    user_id     bigint  primary key,

    disable_levelup_messages    boolean not null    default false,
    reminders_in_dm             boolean not null    default false,
    usernames_opt_out           boolean not null    default false,
    embedless_reminders         boolean not null    default false,
    reaction_pages              boolean not null    default false,

    todo_channel    bigint  not null    default 0
);

create table if not exists todos (
    id  serial  primary key,

    user_id     bigint  not null,
    description text    not null,

    orig_mid        bigint  not null,
    orig_channel_id bigint  not null,
    orig_server_id  bigint  not null,

    mid         bigint  not null,
    channel_id  bigint  not null,
    server_id   bigint  not null,

    complete bool not null default false,

    created     timestamp   not null default (current_timestamp at time zone 'utc'),
    completed   timestamp   default null
);

create table if not exists notes (
    id  serial  primary key,

    server_id   bigint not null,
    user_id     bigint not null,

    note        text        not null,
    moderator   bigint      not null,
    created     timestamp   not null default (current_timestamp at time zone 'utc')
);



create table if not exists channel_mirror (
    server_id       bigint not null,

    from_channel    bigint  primary key,
    to_channel      bigint  not null,

    webhook_id  bigint  not null,
    token       text    not null
);

create table if not exists channel_mirror_messages (
    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  primary key,
    original    bigint  not null,
    user_id     bigint  not null
);

create table if not exists role_categories (
    id          serial  primary key,
    server_id   bigint  not null,

    name        text    not null default '',
    description text    not null default '',
    colour      int     not null default 0,

    require_role    bigint      not null    default 0,
    roles           bigint[]    not null    default array[]::bigint[],

    unique (server_id, name)
);

create table if not exists quotes (
    id  serial  primary key,
    hid char(5) not null unique,

    server_id   bigint  not null,
    channel_id  bigint  not null,
    message_id  bigint  not null,
    user_id     bigint  not null,
    added_by    bigint  not null,
    content     text    not null,
    proxied     boolean not null default false,

    added   timestamp   not null default (current_timestamp at time zone 'utc')
);

create table if not exists highlights (
    user_id     bigint,
    server_id   bigint,
    highlights  text[]      not null    default array[]::text[],
    blocked     bigint[]    not null    default array[]::bigint[],

    primary key (user_id, server_id)
);

create table if not exists highlight_config (
    server_id   bigint      primary key,
    hl_enabled  bool        not null    default false,
    blocked     bigint[]    not null    default array[]::bigint[]
);

create table if not exists highlight_delete_queue (
    message_id  bigint  primary key,
    channel_id  bigint  not null
);

create table if not exists pending_actions (
    id          serial      primary key,
    guild_id    bigint      not null,
    user_id     bigint      not null,
    expires     timestamp   not null,

    -- not an enum for Reasons (we're lazy), can be "unban", "unmute", or "unpause"
    type    text    not null,
    log     bool    not null,
    reason  text    not null
);

-- these messages might have to be deleted 2+ hours after sending
create table if not exists command_responses (
    message_id  bigint  primary key,
    user_id     bigint  not null
);

create table if not exists quote_block (
    user_id bigint  primary key
);

-- +migrate StatementBegin
create or replace function generate_hid() returns char(5) as $$
    select string_agg(substr('abcdefghijklmnopqrstuvwxyz', ceil(random() * 26)::integer, 1), '') from generate_series(1, 5)
$$ language sql volatile;

create or replace function find_free_quote_hid(guild bigint) returns char(5) as $$
declare new_hid char(5);
begin
    loop
        new_hid := generate_hid();
        if not exists (select 1 from quotes where hid = new_hid) then return new_hid; end if;
    end loop;
end
$$ language plpgsql volatile;
-- +migrate StatementEnd

create unique index if not exists tags_uniq_idx on tags (lower(name), server_id);
create index if not exists starboard_messages_id_idx on starboard_messages (starboard_message_id);
create index if not exists usernames_idx on usernames (user_id);
create index if not exists nicknames_idx on nicknames (server_id, user_id);
create index if not exists roles_server_id_idx on roles (server_id);
create index if not exists channel_bans_server_member_idx on channel_bans (server_id, user_id);
create index if not exists reminders_user_idx on reminders (user_id);
create index if not exists reminders_user_server_idx on reminders (user_id, server_id);
create index if not exists nolevels_server_idx on nolevels (server_id);
create index if not exists react_roles_server_idx on react_roles (server_id);
create index if not exists react_role_entries_message_id_idx on react_role_entries (message_id);
create index if not exists todos_user_idx on todos (user_id);
create index if not exists notes_server_idx on notes (server_id);
create index if not exists notes_server_user_idx on notes (server_id, user_id);
create index if not exists role_categories_server_idx on role_categories (server_id);
create index if not exists pending_actions_guild_idx on pending_actions (guild_id);

insert into bot_settings (id, status) values (1, 'online') on conflict (id) do nothing;
