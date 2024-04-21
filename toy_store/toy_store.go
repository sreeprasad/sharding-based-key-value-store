package toy_store

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ToyStoreRecord struct {
	ID        uint
	key       string
	value     string
	expiredAt time.Time
}

type ToyStore struct {
	db *sql.DB
}

func NewToyStore(db *sql.DB) *ToyStore {
	return &ToyStore{db: db}
}

func (toystore *ToyStore) Set(key string, value string, expiredAt time.Time) (bool, error) {

	sqlStatement := `
	INSERT INTO public.toy_dynamo (key, value, expired_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (key) DO UPDATE
	SET value = EXCLUDED.value, expired_at = EXCLUDED.expired_at;
	`
	_, err := toystore.db.Exec(sqlStatement, key, value, expiredAt)
	if err != nil {
		log.Printf("Failed to insert or update record: %v", err)
		return false, err
	}
	return true, nil

}

func (toystore *ToyStore) Get(key string) (ToyStoreRecord, error) {

	sqlStatement := `SELECT id, key, value, expired_at FROM toy_dynamo 
	WHERE key = $1 AND expired_at > NOW()`

	var record ToyStoreRecord
	err := toystore.db.QueryRow(sqlStatement, key).Scan(&record.key, &record.value, &record.expiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No record found for key: %v", key)
			return ToyStoreRecord{}, err
		}
		log.Printf("Failed to retrieve record: %v", err)
		return ToyStoreRecord{}, err
	}
	return record, nil
}
