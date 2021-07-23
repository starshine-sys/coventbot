-- these messages might have to be deleted 2+ hours after sending
create table command_responses (
    message_id  bigint  primary key,
    user_id     bigint  not null
);

---- create above / drop below ----

drop table command_responses;