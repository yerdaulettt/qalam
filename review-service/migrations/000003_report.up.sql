create table if not exists reports (
    id serial primary key,
    problem varchar(300) not null,
    review_id int references reviews(id) on delete cascade,
    reported_at timestamptz default now()
);