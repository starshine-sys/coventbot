alter table ticket_categories add column name text not null default '';

---- create above / drop below ----

alter table ticket_categories drop column name;