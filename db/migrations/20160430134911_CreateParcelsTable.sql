
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table parcels (
  id uuid primary key default gen_random_uuid(),
  dealer_id uuid not null references dealers (id),
  received_date timestamp not null,
  status bool,
  owner_id uuid not null references users (id),
  receiver_id uuid references users (id),
  created_on timestamp,
  modified_on timestamp
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table parcels;

