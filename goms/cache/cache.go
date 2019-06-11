package cache

import (
	"encoding/hex"
	"fmt"
	"hash"
)

type Cache interface {
	Set(key string, value interface{}) (err error)
	Get(key string) (value interface{}, err error)
}

func Key(hasher func() hash.Hash, keys ...interface{}) (key string, err error) {
	h := hasher()
	_, err = fmt.Fprintf(h, "%+v", keys...)
	if err != nil {
		return
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
