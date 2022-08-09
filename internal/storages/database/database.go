package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
)

type DataBaseStorage struct {
	db *sql.DB
	l  *logger.Logger
}

func NewDataBaseStorage(l *logger.Logger, storageConfig config.Storage) *DataBaseStorage {
	st := &DataBaseStorage{l: l}
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(storageConfig.MaxRetries))
	operation := func() error {
		l.Debug("Trying to connect to DB")
		db, err := sql.Open("pgx", storageConfig.DataSourceName)
		if err != nil {
			l.Debug(err)
			return err
		}
		if err = db.Ping(); err != nil {
			l.Debug(err)
			return err
		}
		l.Debug("DB connection is established")
		st.db = db
		return nil
	}
	if err := backoff.Retry(operation, b); err != nil {
		l.Fatal(err)
	}
	return st
}

func (st *DataBaseStorage) Close() (err error) {
	err = st.db.Close()
	if err != nil {
		st.l.Debug(err)
	} else {
		st.l.Debug("Database connection closed correctly")
	}
	return
}

func (st *DataBaseStorage) GetUserURLs(user string) (result models.UserURLs, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := st.db.QueryContext(ctx, getUserURLsQuery, user)
	if err != nil {
		st.l.Debug(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var res models.UserURL
		err = rows.Scan(&res.URL, &res.ShortURL)
		if err != nil {
			continue
		}
		result = append(result, res)
	}
	err = rows.Err()
	if err != nil {
		st.l.Debug(err)
	}
	return
}

type webResourse struct {
	webResourseID string
	URL           string
	shortURL      string
	counter       int
	userID        string
	isDeleted     bool
}

func (st *DataBaseStorage) GetByShortURL(requiredURL string) (URL string, err error) {
	result := webResourse{}
	err = st.db.QueryRow(getURLQuery, requiredURL).Scan(
		&result.webResourseID,
		&result.URL,
		&result.shortURL,
		&result.counter,
		&result.userID,
		&result.isDeleted)
	if err != nil {
		err = errors.New(requiredURL + " does not exist")
		return
	}
	st.l.Debug(result, err)
	URL = result.URL
	_, err = st.db.Exec(updateCounterQuery, result.counter+1, result.webResourseID)
	if err != nil {
		st.l.Debug(err)
	}
	err = nil
	if result.isDeleted {
		err = shortenererrors.ErrURLDeleted
	}
	return
}

func (st *DataBaseStorage) saveUser(userID string) {
	if !checkValueExists(st.db, "user", "user_id", userID) {
		_, err := st.db.Exec(insertUserQuery, userID)
		if err != nil {
			st.l.Debug(err)
		}
	}
}

func (st *DataBaseStorage) SaveURL(URL string, shortURL string, userID string) (err error) {
	st.saveUser(userID)
	webResourseUUID, err := uuid.NewV4()
	if err != nil {
		st.l.Debug(err)
	}
	_, err = st.db.Exec(insertWebResourseQuery,
		webResourseUUID.String(), URL, shortURL, "0", userID)

	if err == nil {
		return
	}
	st.l.Debug(err)
	if pqerr, ok := err.(*pgconn.PgError); ok {
		if pgerrcode.IsIntegrityConstraintViolation(pqerr.Code) {
			return storageerrors.ErrURLAlreadySave
		}
	}
	return
}

func (st *DataBaseStorage) SaveBatch(
	batch models.BatchURLs,
	userID string) (result models.BatchShortURLs, err error) {
	st.saveUser(userID)

	result = make(models.BatchShortURLs, len(batch))
	tx, err := st.db.Begin()
	if err != nil {
		st.l.Debug(err)
		return
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		st.l.Debug(err)
	}(tx)

	stmt, err := tx.PrepareContext(context.Background(), insertWebResourseBatchQuery)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			st.l.Debug("Close statement error")
		}
	}(stmt)

	for index, row := range batch {
		webResourseUUID, err := uuid.NewV4()
		if err != nil {
			st.l.Debug(err)
		}
		result[index].ShortURL = row.ShortURL
		result[index].CorrelationID = row.CorrelationID

		if _, err = stmt.ExecContext(
			context.Background(),
			webResourseUUID.String(),
			row.OriginalURL,
			row.ShortURL,
			0,
			userID); err != nil {
			st.l.Debug(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		st.l.Debug(err)
	}
	return
}

func (st *DataBaseStorage) CheckConnection(ctx context.Context) error {
	return st.db.PingContext(ctx)
}

func (st *DataBaseStorage) DeleteBatch(ctx context.Context, batch models.BatchForDelete, user string) (err error) {
	tx, err := st.db.Begin()
	if err != nil {
		st.l.Debug(err)
		return
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		st.l.Debug(err)
	}(tx)
	stmt, err := tx.PrepareContext(ctx, batchUpdate)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			st.l.Debug("Close statement error")
		}
	}(stmt)
	if _, err = stmt.ExecContext(
		ctx, pq.Array(batch), user); err != nil {
		st.l.Debug(err)
	}
	err = tx.Commit()
	st.l.Debug("deleted from DB", batch)
	return
}

func checkValueExists(db *sql.DB, table string, field string, value string) bool {
	var count int
	query := fmt.Sprintf(checkValueExistsQuery, field, table, field)
	err := db.QueryRow(query, value).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return true
	}
	return false
}
