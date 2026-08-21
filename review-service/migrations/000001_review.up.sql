create table if not exists books_id (
    id int primary key
);

create table if not exists reviews (
    id serial primary key,
    content text not null,
    user_id int not null,
    book_id int references books_id(id) on delete cascade
);