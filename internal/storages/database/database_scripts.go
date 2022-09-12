package database

const insertUserQuery = `
	INSERT INTO
 		"public"."user" (user_id)
	VALUES ($1);`

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
    ($1, $2, $3, $4, $5);
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
		%s=$1;
`

const getURLQuery = `
SELECT
    web_resourse_id,
    url,
    short_url,
    counter,
    user_id,
    is_deleted
FROM
    public.web_resourse
WHERE
    short_url = $1;
`
const getUserURLsQuery = `
SELECT
    url,
    short_url
FROM
    public.web_resourse
WHERE
    user_id = $1;
`
const updateCounterQuery = `
UPDATE
    public.web_resourse
SET
    counter = $1
WHERE
    web_resourse_id = $2;
`

const batchUpdate = `
UPDATE public.web_resourse
SET is_deleted = TRUE
WHERE short_url = any ($1) AND user_id=$2;
`

const countUsers = `
SELECT count(user_id) from public.user;
`

const countShortURLs = `
SELECT count(web_resourse_id) from public.web_resourse;
`
