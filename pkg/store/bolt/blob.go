package bolt

import (
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/misc"
	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/boltdb/bolt"
	"github.com/renstrom/shortuuid"
)

func (bs *store) SaveBlob(raw types.Blob) ([]byte, error) {

	key := []byte(shortuuid.New())
	err := bs.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket(BlobBucket)
		gzipped, e := misc.Zip(raw)
		if e != nil {
			return e
		}

		return bu.Put(key, gzipped)
	})
	return key, err
}

func (b *store) Blob(key []byte) (types.Blob, error) {
	var (
		raw []byte
		err error
	)
	err = b.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket(BlobBucket)
		gzipped := bu.Get(key)
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

func (bs *store) DeleteBlob(key []byte) error {

	// Delete the key in a different write transaction.
	err := bs.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(BlobBucket).Delete(key)
	})
	return err

}
