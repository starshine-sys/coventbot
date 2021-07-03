alter table user_config add column embedless_reminders boolean not null default false;

---- create above / drop below ----

alter table user_config drop column embedless_reminders;