alter table user_config add column reaction_pages boolean not null default false;

---- create above / drop below ----

alter table user_config drop column reaction_pages;