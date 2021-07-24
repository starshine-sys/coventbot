alter table react_roles add column title text;
alter table react_roles add column description text;
alter table react_roles add column mention boolean;

---- create above / drop below ----

alter table react_roles drop column title;
alter table react_roles drop column description;
alter table react_roles drop column mention;