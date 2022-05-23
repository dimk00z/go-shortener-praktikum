package database

const createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS "public"."user" (
    user_id uuid NOT NULL,
    CONSTRAINT idx_table PRIMARY KEY (user_id)
);
`

const createWebResourseTableQuery = `
CREATE TABLE IF NOT EXISTS "public"."web_resourse" (
    web_resourse_id uuid NOT NULL,
    url varchar(300) NOT NULL,
    short_url varchar(50) NOT NULL,
    counter integer,
    user_id uuid NOT NULL,
    CONSTRAINT web_resourse_pkey PRIMARY KEY (web_resourse_id),
    CONSTRAINT user_id_url UNIQUE (user_id, url),
    CONSTRAINT fk_web_resourse_user FOREIGN KEY (user_id) REFERENCES public."user" (user_id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);`

const insertUserQuery = `
	INSERT INTO
 		"public"."user" (user_id)
	VALUES ('%s');`

const insertWebResourseQuery = `
INSERT INTO
    public.web_resourse(
        web_resourse_id,
        url,
        short_url,
        counter,
        user_id
    )
VALUES
    ('%s', '%s','%s', %s, '%s');
`
const insertWebResourseBatchQuery = `
INSERT INTO
    public.web_resourse(
        web_resourse_id,
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
		%s='%s';
`

const getURLQuery = `
SELECT
    web_resourse_id,
    url,
    short_url,
    counter,
    user_id
FROM
    public.web_resourse
WHERE
    short_url = '%s';
`
const getUserURLsQuery = `
SELECT
    url,
    short_url
FROM
    public.web_resourse
WHERE
    user_id = '%s';
`
const updateCounterQuery = `
UPDATE
    public.web_resourse
SET
    counter = %v
WHERE
    web_resourse_id = '%s';
`
