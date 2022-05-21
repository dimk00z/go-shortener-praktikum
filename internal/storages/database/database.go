package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/qustavo/dotsql"
)

type DataBaseStorage struct {
	db          *sql.DB
	sqlFilePath string
}

func NewDataBaseStorage(DataSourceName string, sqlFilePath string) *DataBaseStorage {
	db, err := sql.Open("pgx", DataSourceName)
	if err != nil {
		log.Println(err)
	}
	dot, err := dotsql.LoadFromFile(sqlFilePath)
	if err != nil {
		log.Println(err)
	}

	if _, err := dot.Exec(db, "create-users-table"); err != nil {
		log.Println(err)
	}

	if _, err = dot.Exec(db, "create-web-resourse-table"); err != nil {
		log.Println(err)
	}

	return &DataBaseStorage{
		db:          db,
		sqlFilePath: sqlFilePath,
	}
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
	return
}
func (st *DataBaseStorage) GetByShortURL(requiredURL string) (shortURL string, err error) {
	return
}

func (st *DataBaseStorage) SaveURL(URL string, shortURL string, userID string) {

}

func (st *DataBaseStorage) CheckConnection(ctx context.Context) error {
	return st.db.PingContext(ctx)
}
