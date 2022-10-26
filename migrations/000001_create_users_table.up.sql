CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    users_id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    users_name text NOT NULL,
    users_description text NOT NULL,
    nationality varchar(32) NOT NULL,
    birthday DATE NOT NULL, 
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    verified_pro bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX idx_nationality ON users(nationality);