package urlstore

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const DATABASE_FILE_NAME = "urlstore.db"
const URL_BUCKET = "urls"

type Store interface {
	Get(path string) (*Redirect, error)
	Put(redirect *Redirect) error
}

type UrlStore struct{}

func NewUrlStore() *UrlStore {
	return &UrlStore{}
}

func (*UrlStore) Get(path string) (*Redirect, error) {

	connection, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	var response Redirect

	err = connection.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(URL_BUCKET))

		if bucket == nil {
			//nothing to do here, there are no keys.
			return nil
		}

		val := bucket.Get([]byte(path))

		if val == nil {
			return fmt.Errorf("unable to find key '%s' in url store", path)
		}

		response = Redirect{
			path,
			string(val),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (*UrlStore) Put(redirect *Redirect) error {

	connection, err := getConnection()
	if err != nil {
		return err
	}
	defer connection.Close()

	err = connection.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(URL_BUCKET))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(redirect.Path), []byte(redirect.Url))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getConnection() (*bolt.DB, error) {
	return bolt.Open(DATABASE_FILE_NAME, 0600, nil)
}
