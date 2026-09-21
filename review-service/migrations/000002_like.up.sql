create table if not exists review_likes (
    review_id int references reviews(id),
    user_id int references users(id),
    liked boolean not null default false,
    primary key(review_id, user_id)
);