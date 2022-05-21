package database

const createUsersTable = `
	CREATE TABLE IF NOT EXISTS "public"."user" (
    user_id uuid NOT NULL,
    CONSTRAINT idx_table PRIMARY KEY (user_id)
);`

const createWebResourseTable = `
	CREATE  TABLE "public".web_resourse ( 
	web_resourse_id      uuid  NOT NULL  ,
	url                  varchar(300) NOT NULL ,
	short_url            varchar(50) NOT NULL ,
	counter              integer    ,
	user_id              uuid    ,
	CONSTRAINT idx_web_resourse PRIMARY KEY ( web_resourse_id ),
	CONSTRAINT fk_web_resourse_user FOREIGN KEY ( user_id ) REFERENCES "public"."user"( user_id ) ON DELETE CASCADE ON UPDATE CASCADE 
 );`
