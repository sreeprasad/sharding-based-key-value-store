package toy_store

import (
	"database/sql"
	"hash/fnv"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ToyStoreRecord struct {
	ID        uint
	Key       string
	Value     string
	ExpiredAt time.Time
}

type ToyStore struct {
	db1 *sql.DB
	db2 *sql.DB
}

func NewToyStore(db1 *sql.DB, db2 *sql.DB) *ToyStore {
	return &ToyStore{db1: db1, db2: db2}
}

func (toystore *ToyStore) Set(key string, value string, expiredAt time.Time) (bool, error) {

	sqlStatement := `
	INSERT INTO public.toy_dynamo (key, value, expired_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (key) DO UPDATE
	SET value = EXCLUDED.value, expired_at = EXCLUDED.expired_at;
	`
	db := getShard(toystore, key)
	_, err := db.Exec(sqlStatement, key, value, expiredAt)
	if err != nil {
		log.Printf("Failed to insert or update record: %v", err)
		return false, err
	}
	return true, nil

}

func getShard(toystore *ToyStore, key string) *sql.DB {

	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	hashValue := hasher.Sum32()

	if hashValue%2 == 0 {
		log.Printf("connecting to shard 1 for key: %s", key)
		return toystore.db1
	} else {
		log.Printf("connecting to shard 2 for key: %s", key)
		return toystore.db2
	}
}

func (toystore *ToyStore) Get(key string) (ToyStoreRecord, error) {

	sqlStatement := `SELECT id, key, value, expired_at FROM toy_dynamo 
	WHERE key = $1 AND expired_at > NOW()`

	var record ToyStoreRecord
	db := getShard(toystore, key)
	err := db.QueryRow(sqlStatement, key).Scan(&record.ID, &record.Key, &record.Value, &record.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No record found for key: %v", key)
			return ToyStoreRecord{}, nil
		}
		log.Printf("Failed to retrieve record: %v", err)
		return ToyStoreRecord{}, err
	}
	return record, nil
}

func (toystore *ToyStore) Delete(key string) (bool, error) {

	sqlStatement := `
	UPDATE toy_dynamo set expired_at = TIMESTAMP '1970-01-01 00:00:00' where key = $1 
	and expired_at > NOW();
	`
	db := getShard(toystore, key)
	_, err := db.Exec(sqlStatement, key)
	if err != nil {
		log.Printf("Failed to delete record for key:%s due to error: %v", key, err)
		return false, err
	}
	return true, nil

}
