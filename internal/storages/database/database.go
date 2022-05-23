package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/gofrs/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DataBaseStorage struct {
	db *sql.DB
}

func NewDataBaseStorage(dbConfig settings.DBStorageConfig) *DataBaseStorage {
	st := &DataBaseStorage{}
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(dbConfig.MaxRetries))
	operation := func() error {
		log.Println("Trying to connect to DB")
		db, err := sql.Open("pgx", dbConfig.DataSourceName)
		if err != nil {
			log.Println(err)
			return err
		}
		if err = db.Ping(); err != nil {
			log.Println(err)
			return err
		}
		log.Println("DB connection is established")
		st.db = db
		return nil
	}
	if err := backoff.Retry(operation, b); err != nil {
		log.Panicln(err)
	}
	createTables(st.db, createUsersTableQuery, createWebResourseTableQuery)
	return st
}

func (st *DataBaseStorage) Close() (err error) {
	err = st.db.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database connection closed correctly")
	}
	return
}

func (st *DataBaseStorage) GetUserURLs(user string) (result []struct {
	ShortURL string
	URL      string
}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := st.db.QueryContext(ctx, fmt.Sprintf(getUserURLsQuery, user))
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var res struct {
			ShortURL string
			URL      string
		}
		err = rows.Scan(&res.URL, &res.ShortURL)
		if err != nil {
			continue
		}
		result = append(result, res)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

type webResourse struct {
	webResourseID string
	URL           string
	shortURL      string
	counter       int
	userID        string
}

func (st *DataBaseStorage) GetByShortURL(requiredURL string) (URL string, err error) {
	result := webResourse{}
	err = st.db.QueryRow(fmt.Sprintf(getURLQuery, requiredURL)).Scan(
		&result.webResourseID, &result.URL, &result.shortURL, &result.counter, &result.userID)
	if err != nil {
		err = errors.New(requiredURL + " does not exist")
		return
	}
	log.Println(result, err)
	URL = result.URL
	_, err = st.db.Exec(fmt.Sprintf(updateCounterQuery, result.counter+1, result.webResourseID))
	if err != nil {
		log.Println(err)
	}
	err = nil
	return
}

func (st *DataBaseStorage) SaveURL(URL string, shortURL string, userID string) {
	// addUser
	if !checkValueExists(st.db, "user", "user_id", userID) {
		_, err := st.db.Exec(fmt.Sprintf(insertUserQuery, userID))
		if err != nil {
			log.Println(err)
		}
	}
	webResourseUUID, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	_, err = st.db.Exec(
		fmt.Sprintf(insertWebResourseQuery,
			webResourseUUID.String(), URL, shortURL, "0", userID))
	if err != nil {
		log.Println(err)
	}

}

func (st *DataBaseStorage) SaveBatch(
	batch models.BatchURLs,
	user string) (result models.BatchShortURLs, err error) {
	result = make(models.BatchShortURLs, len(batch))
	tx, err := st.db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	// var insertStmt *sql.Stmt
	// TODO add batch
	defer tx.Rollback()

	return result, err
}

func (st *DataBaseStorage) CheckConnection(ctx context.Context) error {
	return st.db.PingContext(ctx)
}

func createTables(db *sql.DB, tables ...string) {
	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			log.Println(err)
		}
	}
}

func checkValueExists(db *sql.DB, table string, field string, value string) bool {
	var count int
	err := db.QueryRow(fmt.Sprintf(checkValueExistsQuery, field, table, field, value)).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return true
	}
	return false
}
