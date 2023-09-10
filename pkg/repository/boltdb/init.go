package boltdb

import (
	"github.com/LittleMikle/TG_bot/pkg/repository"
	"github.com/boltdb/bolt"
)

func ConnectToBoltDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
