//Package bolt 实现基于blotdb的存储设备
package bolt

import (
	"os"
	"time"

	"github.com/lewgun/argyroneta/pkg/store"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	//"github.com/pquerna/ffjson/ffjson"
)

var (
    
    //RuleBucket is a rule bucket
	RuleBucket  = []byte("rule")
    
    //BlobBucket is a blob bucket
	BlobBucket  = []byte("blob")
)

//store a Store implemented with blot as backend
type store struct {
	db     *bolt.DB
	opened bool
}

//Connect open a bolt instance
func (bs *store) Connect(filePath string) error {
	var err error
	bs.db, err = bolt.Open(
		filePath,
		0600,
		&bolt.Options{
			Timeout: 1 * time.Second,
		})

	if err != nil {
		return err
	}

	err = bs.db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(RuleBucket)
		if e != nil {
			return e
		}

		_, e = tx.CreateBucketIfNotExists(BlobBucket)
		if e != nil {
			return e
		}
		return nil
	})

	go func() {
		// Grab the initial stats.
		prev := bs.db.Stats()

		for {

			// Grab the current stats and diff them.
			stats := bs.db.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
		//	ffjson.NewEncoder(os.Stderr).Encode(diff)
			// Save stats for the next loop.
			prev = stats

			// Wait for 10s.
			time.Sleep(60 * time.Second)
		}
	}()

	bs.opened = true

	return
}

//Close close the blot instance
func (bs *store) Close() error {
	bs.db.Close()
	bs.opened = false
}

func init() {
	store.Register(store.Bolt, &store{})
}
