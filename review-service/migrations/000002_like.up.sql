create table if not exists review_likes (
    review_id int references reviews(id) on delete cascade,
    user_id int references users(id) on delete cascade,
    liked boolean not null default false,
    primary key(review_id, user_id)
);