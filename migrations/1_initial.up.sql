create schema if not exists tasks_app;

create type task_state AS ENUM ('created', 'processing', 'accepted', 'rejected');


create table if not exists tasks_app.tasks_state
(
    id              varchar(24) primary key not null,
    state           task_state default 'created',
    time_agreement  interval
);

create table if not exists tasks_app.user_accept_time
(
    id              serial primary key,
    task_id         varchar(24) not null references tasks_app.tasks_state(id),
    email           varchar(100) not null,
    time_start      timestamp,
    time_end        timestamp
);