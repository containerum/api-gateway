CREATE TABLE IF NOT EXISTS "groups" (
  id uuid DEFAULT uuid_generate_v1() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
  updated_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
  name varchar(32) NOT NULL,
  active boolean NOT NULL DEFAULT false
);
