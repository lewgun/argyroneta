package bolt

import (
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/misc"

	"github.com/boltdb/bolt"
)

func (bs *store) SaveBlob(key string, raw []byte) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket(BlobBucket)
		gzipped, e := misc.Zip(raw)
		if e != nil {
			return e
		}
		return bu.Put([]byte(key), gzipped)
	})
	return err
}

func (b *store) Blob(key string) ([]byte, error) {
	var (
		raw []byte
		err error
	)
	err = b.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket(BlobBucket)
		gzipped := bu.Get([]byte(key))
		if len(gzipped) == 0 {
			return errutil.ErrNotFound
		}
		raw, err = misc.Unzip(gzipped)
		if err != nil {
			return err
		}
		return nil
	})

	return raw, err

}


func (bs *store) DeleteBlob(key string) error {

	// Delete the key in a different write transaction.
	err := bs.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(RuleBucket).Delete([]byte(key))
	})
	return err

}
