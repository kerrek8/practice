

create table if not exists users (
    id INTEGER primary key,
    username text not null unique,
    password VARCHAR(256) not null,
    name text not null
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
