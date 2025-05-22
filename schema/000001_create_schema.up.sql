CREATE TABLE users (
    id serial primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists (
    id serial primary key,
    title varchar(255) not null,
    description varchar(255)
);

CREATE TABLE user_lists (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    list_id int not null references todo_lists(id) on delete cascade
);

CREATE TABLE todo_items (
    id serial primary key,
    title varchar(255) not null,
    description varchar(255),
    done boolean not null default false
);

CREATE TABLE lists_items (
    id serial primary key,
    item_id int not null references todo_items(id) on delete cascade,
    list_id int not null references todo_lists(id) on delete cascade
);
