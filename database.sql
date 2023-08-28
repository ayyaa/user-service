/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
-- CREATE TABLE test (
-- 	id serial PRIMARY KEY,
-- 	name VARCHAR ( 50 ) UNIQUE NOT NULL,
-- );

-- INSERT INTO test (name) VALUES ('test1');
-- INSERT INTO test (name) VALUES ('test2');

CREATE TABLE public.users
(
    id bigserial NOT NULL,
    uuid character varying(60),
    slug character varying(255),
    email character varying(100),
    password character varying(60),
    name character varying(60),
    phone character varying(13),
    updated_at time without time zone,
    created_at time without time zone,
    CONSTRAINT pk_users_id PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;


INSERT INTO users (uuid, slug, email, name, password, name, phone, updated_at, created_at) 
VALUES (uuid_generate_v4(), 'ayya', 'tsurayya@gmail.com', 'tsurayya', '123', '+6281219823417', NOW(), NOW());
