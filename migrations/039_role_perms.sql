alter table servers add column helper_roles bigint[] not null default array[]::bigint[];
alter table servers add column mod_roles bigint[] not null default array[]::bigint[];
alter table servers add column admin_roles bigint[] not null default array[]::bigint[];

alter table servers add column roles_set_up boolean not null default false;

---- create above / drop below ----

alter table servers drop column helper_roles;
alter table servers drop column mod_roles;
alter table servers drop column admin_roles;
alter table servers drop column roles_set_up;