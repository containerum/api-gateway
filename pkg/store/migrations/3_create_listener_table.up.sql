CREATE TYPE "method" AS ENUM ('GET', 'POST', 'PUT', 'DELETE', 'PATCH');


CREATE TABLE IF NOT EXISTS "listeners" (
  id uuid DEFAULT uuid_generate_v1() PRIMARY KEY,
  created_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
  updated_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
  name varchar(128) NOT NULL,
  o_auth boolean NOT NULL DEFAULT true,
  active boolean NOT NULL DEFAULT false,
  group_refer uuid REFERENCES groups(id),
  strip_path boolean NOT NULL DEFAULT false,
  listen_path varchar(128) NOT NULL,
  upstream_url varchar(128) NOT NULL,
  method method NOT NULL
);
