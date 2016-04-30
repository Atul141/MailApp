
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table dealers (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  icon text,
  created_on timestamp,
  modified_on timestamp
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table dealers;
