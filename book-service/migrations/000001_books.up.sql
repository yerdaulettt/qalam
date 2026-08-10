create table if not exists authors (
    id serial primary key,
    name varchar(25) not null,
    surname varchar(25) not null,
    about text not null
);

create table if not exists books (
    id serial primary key,
    name varchar(150) not null,
    description text not null,
    author_id int references authors(id) on delete cascade
);

create table if not exists genres (
    id serial primary key,
    name varchar(30) not null unique
);

create table if not exists book_genres (
    book_id int references books(id) on delete cascade,
    genre_id int references genres(id) on delete cascade,
    primary key(book_id, genre_id)
);