

create table if not exists users (
    id INTEGER primary key,
    username text not null unique,
    password VARCHAR(256) not null,
    name text not null,
    role TEXT DEFAULT 'agent'
);

create table if not exists listings (
    id INTEGER primary key,
    name text not null,
    type text not null,
    description text not null,
    status text not null,
    price numeric not null,
    city text not null,
    user_id integer,
    date_created datetime not null default (datetime('now', '+5 hours')),
    foreign key (user_id) references users(id) on delete cascade
);


INSERT INTO users (username, password, name, role) VALUES ('admin', '$2a$10$Pl3VNUeA.4vyh1ZJQvUC1O18damR8YZhCHEW0jF1icQexSLJgCn62', 'Administrator', 'admin');