create table if not exists users (
    id serial primary key,
    name varchar(20) not null check (name != ''),
    username varchar(30) unique not null check (username != ''),
    role varchar(20) constraint role_check check (role in ('user', 'admin')),
    hash text not null check (hash != '')
);