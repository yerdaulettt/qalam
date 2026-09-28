create table if not exists book_ratings (
    book_id int references books(id),
    user_id int not null,
    rating int check (rating > 0 and rating <= 10),
    primary key(book_id, user_id)
);