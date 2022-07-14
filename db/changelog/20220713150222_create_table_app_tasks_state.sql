-- +goose Up

create type task_state AS ENUM ('created', 'processing', 'accepted', 'rejected');

create table tasks_app.tasks_state
(
    id varchar(24) primary key not null,
    state task_state default 'created',
    time_agreement interval
);