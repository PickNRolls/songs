create table songs (
  id serial primary key,  
  group_name text not null,
  name text not null,
  verses text[] not null,
  link text,
  released_at timestamp not null,
  added_at timestamp not null,
  updated_at timestamp not null
);

