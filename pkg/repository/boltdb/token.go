package boltdb

import (
	"errors"
	"github.com/LittleMikle/TG_bot/pkg/repository"
	"github.com/boltdb/bolt"
	"strconv"
)

type TokenRepoBolt struct {
	db *bolt.DB
}

func NewTokenRepoBolt(db *bolt.DB) *TokenRepoBolt {
	return &TokenRepoBolt{
		db: db,
	}
}

func (r *TokenRepoBolt) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})
}

func (r *TokenRepoBolt) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
