
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION pgcrypto;
create table users (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  email text not null,
  phone_no text,
  emp_id text not null,
  created_on timestamp,
  modified_on timestamp
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop EXTENSION pgcrypto;
drop table users;
