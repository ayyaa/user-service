/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE public.users
(
    id bigserial NOT NULL,
    password character varying(60),
    full_name character varying(60),
    phone_number character varying(15) UNIQUE,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_users_id PRIMARY KEY (id),
    CONSTRAINT uq_users_phone UNIQUE (phone_number)
);
