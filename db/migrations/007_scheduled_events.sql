-- +migrate Up

-- 2022-01-24
-- Add scheduled events
-- Also remove the old pending actions table

create table scheduled_events (
    id          serial      primary key,
    event_type  text        not null,
    expires     timestamp   not null,
    data        jsonb       not null
);

create index scheduled_events_expires_idx on scheduled_events (expires);
create index scheduled_events_data_idx on scheduled_events using GIN (data);

drop table pending_actions;
