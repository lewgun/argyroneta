package bolt

import (
	"sync"

	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/Sirupsen/logrus"
)


type BlobStore interface {
	SaveBlob(key []byte, blob types.Blob) error
	DeleteBlob(key []byte) error
	Blob(key []byte) (types.Blob, error)
}

//Store 为所有存储设备提供一个基本接口
type Store interface {
	BlobStore

	PowerOff() error
}


//New new a blotdb instance
func New( conf string, logger *logrus.Logger) Store {
	s := &store{}
	
	if err := s.init(conf, logger); err != nil {
		panic(err)
	}
	return s
}



//store a Store implemented with blot as backend
type store struct {
	db     *bolt.DB
	logger *logrus.Logger
}

//PowerOn open a bolt instance
func (bs *store) init(filePath string, logger *logrus.Logger) error {

	bs.logger = logger

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
	//
	//go func() {
	//	// Grab the initial stats.
	//	prev := bs.db.Stats()
	//
	//	for {
	//
	//		// Grab the current stats and diff them.
	//		stats := bs.db.Stats()
	//		diff := stats.Sub(&prev)
	//
	//		// Encode stats to JSON and print to STDERR.
	//		//	ffjson.NewEncoder(os.Stderr).Encode(diff)
	//		// Save stats for the next loop.
	//		prev = stats
	//
	//		// Wait for 10s.
	//		time.Sleep(60 * time.Second)
	//	}
	//}()

	return nil
}

//Close close the blot instance
func (bs *store) Close() error {
	return bs.db.Close()

}

func 