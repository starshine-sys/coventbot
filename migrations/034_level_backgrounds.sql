alter table levels add column background text not null default '';
alter table server_levels add column background text not null default '';

---- create above / drop below ----

alter table levels drop column background;
alter table server_levels drop column background;