-- +goose Up

create table tasks_app.user_accept_time
(
    id         serial primary key,
    task_id    varchar(24) not null references tasks_app.tasks_state(id),
    email      varchar(100) not null,
    time_start      timestamp,
    time_end        timestamp
);