create table if not exists "movies" (
    id integer primary key autoincrement,
    path text not null,
    title text not null,
    year integer not null,
    library_id integer,
    foreign key (library_id) references "libraries" (id) on delete cascade on update cascade
);
