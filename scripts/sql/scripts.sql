-- name: create-users-table
CREATE TABLE IF NOT EXISTS "public"."user" (
    user_id uuid NOT NULL,
    CONSTRAINT idx_table PRIMARY KEY (user_id)
);

-- name: create-web-resourse-table
CREATE TABLE IF NOT EXISTS "public"."web_resourse" (
    web_resourse_id uuid NOT NULL,
    url varchar(300),
    short_url varchar(50),
    counter integer,
    user_id uuid,
    CONSTRAINT idx_web_resourse PRIMARY KEY (web_resourse_id),
    CONSTRAINT url_short_url_unique UNIQUE (url, short_url),
    CONSTRAINT fk_web_resourse_user FOREIGN KEY (user_id) REFERENCES "public"."user"(user_id)
);