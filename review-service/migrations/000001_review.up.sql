create table if not exists books (
    id int primary key,
    name varchar(150) not null
);

create table if not exists users (
    id int primary key,
    username varchar(30) unique not null
);

create table if not exists reviews (
    id serial primary key,
    content text not null,
    user_id int references users(id) on delete cascade,
    book_id int references books(id) on delete cascade,
);