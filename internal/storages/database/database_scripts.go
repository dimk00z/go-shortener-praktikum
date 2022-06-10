package database

const createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS "public"."user" (
    user_id uuid NOT NULL,
    CONSTRAINT idx_table PRIMARY KEY (user_id)
);
`

const createWebResourceTableQuery = `
CREATE TABLE IF NOT EXISTS "public"."web_Resource" (
    web_Resource_id uuid NOT NULL,
    url varchar(300) NOT NULL,
    short_url varchar(50) NOT NULL,
    counter integer,
    user_id uuid NOT NULL,
    CONSTRAINT web_Resource_pkey PRIMARY KEY (web_Resource_id),
    CONSTRAINT user_id_url UNIQUE (user_id, url),
    CONSTRAINT fk_web_Resource_user FOREIGN KEY (user_id) REFERENCES public."user" (user_id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);`

const insertUserQuery = `
	INSERT INTO
 		"public"."user" (user_id)
	VALUES ($1);`

const insertWebResourceQuery = `
INSERT INTO
    public.web_Resource(
        web_Resource_id,
        url,
        short_url,
        counter,
        user_id
    )
VALUES
    ($1, $2, $3, $4, $5);
`
const insertWebResourceBatchQuery = `
INSERT INTO
    public.web_Resource(
        web_Resource_id,
        url,
        short_url,
        counter,
        user_id
    )
VALUES
    ($1, $2,$3, $4, $5)
	ON CONFLICT (user_id, url) DO NOTHING;`

const checkValueExistsQuery = `
	SELECT 
		COUNT(%s)
	FROM 
		"public"."%s"
	WHERE
		%s=$1;
`

const getURLQuery = `
SELECT
    web_Resource_id,
    url,
    short_url,
    counter,
    user_id
FROM
    public.web_Resource
WHERE
    short_url = $1;
`
const getUserURLsQuery = `
SELECT
    url,
    short_url
FROM
    public.web_Resource
WHERE
    user_id = $1;
`
const updateCounterQuery = `
UPDATE
    public.web_Resource
SET
    counter = $1
WHERE
    web_Resource_id = $2;
`
