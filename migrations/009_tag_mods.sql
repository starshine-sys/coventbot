alter table servers add column tag_mod_role bigint not null default 0;

---- create above / drop below ----

alter table servers drop column tag_mod_role;