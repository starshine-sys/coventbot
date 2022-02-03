-- +migrate Up

-- 2022-02-03
-- Add user, guild, user+guild stores
-- NOTE: if you're running a PostgreSQL version pre-13, and the user you're connecting with is not a superuser, you must run `CREATE EXTENSION hstore;` on the database manually *BEFORE* running this migration.

-- +migrate StatementBegin
DO $$
DECLARE
version integer := (regexp_match(version(), 'PostgreSQL (\d+)\.*'))[1];
superuser boolean := (select usesuper from pg_user where usename = CURRENT_USER);
BEGIN
IF coalesce(version >= 13, true) OR superuser THEN
    create extension if not exists hstore;
END IF;
END $$;
-- +migrate StatementEnd

create table user_config_new (
    user_id bigint  primary key,
    config  hstore  not null default ''
);

-- +migrate StatementBegin
DO $$
DECLARE
row user_config%rowtype;
BEGIN
    FOR row in SELECT * FROM user_config
	LOOP
        insert into user_config_new (user_id, config)
        values (row.user_id, hstore(row));
    END LOOP;
END $$;
-- +migrate StatementEnd

update user_config_new set config = delete(config, 'user_id');

alter table user_config rename to user_config_old;
alter table user_config_new rename to user_config;

create table guild_config (
    guild_id    bigint  primary key,
    config      hstore  not null default ''
);

create table user_guild_config (
    user_id     bigint  not null,
    guild_id    bigint  not null,
    config      hstore  not null default '',

    primary key (guild_id, user_id)
);
