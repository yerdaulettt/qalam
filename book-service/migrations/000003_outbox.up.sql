create table if not exists outbox (
    id uuid primary key default gen_random_uuid(),
    event_type varchar(20) not null,
    event jsonb not null
);