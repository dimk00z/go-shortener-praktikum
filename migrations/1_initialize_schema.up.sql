CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "public"."user" (
    user_id uuid NOT NULL,
    CONSTRAINT idx_table PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS "public"."web_resourse" (
    web_resourse_id uuid NOT NULL,
    url varchar(300) NOT NULL,
    short_url varchar(50) NOT NULL,
    counter integer,
    is_deleted boolean DEFAULT FALSE,
    user_id uuid NOT NULL,
    CONSTRAINT web_resourse_pkey PRIMARY KEY (web_resourse_id),
    CONSTRAINT user_id_url UNIQUE (user_id, url),
    CONSTRAINT fk_web_resourse_user FOREIGN KEY (user_id) REFERENCES public."user" (user_id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);