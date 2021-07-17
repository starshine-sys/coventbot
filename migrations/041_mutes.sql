create table pending_actions (
    id          serial      primary key,
    guild_id    bigint      not null,
    user_id     bigint      not null,
    expires     timestamp   not null,

    -- not an enum for Reasons (we're lazy), can be "unban", "unmute", or "unpause"
    type    text    not null,
    log     bool    not null,
    reason  text    not null
);

create index pending_actions_guild_idx on pending_actions (guild_id);

alter table servers add column mute_role bigint not null default 0;
alter table servers add column pause_role bigint not null default 0;
alter table servers add column muteme_message text not null default 'Successfully {action}d {mention} for {duration}.'

---- create above / drop below ----

drop table pending_actions;
alter table servers drop column mute_role;
alter table servers drop column pause_role;
alter table servers drop column muteme_message;