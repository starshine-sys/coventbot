alter table levels add column colour bigint not null default 0;

---- create above / drop below ----

alter table levels drop column colour;