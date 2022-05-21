package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DataBaseStorage struct {
	db *sql.DB
}

func NewDataBaseStorage(DataSourceName string) *DataBaseStorage {
	db, err := sql.Open("pgx", DataSourceName)
	if err != nil {
		log.Println(err)
	}
	createTables(db, createUsersTable, createWebResourseTable)

	return &DataBaseStorage{
		db: db,
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

func createTables(db *sql.DB, tables ...string) {
	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			log.Println(err)
		}
	}
}
