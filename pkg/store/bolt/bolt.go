package bolt

import (
	"time"

	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

var (
	BlobBucket = []byte("blob")
)

type BlobStore interface {
	SaveBlob(blob types.Blob) ([]byte, error)
	DeleteBlob(key []byte) error
	Blob(key []byte) (types.Blob, error)
}

//Store 为所有存储设备提供一个基本接口
type Store interface {
	BlobStore

	PowerOff() error
}

//New new a blotdb instance
func New(conf *types.BoltConf, logger *logrus.Logger) Store {

	if conf == nil || logger == nil {
		return nil
	}
	s := &store{}

	if err := s.init(conf, logger); err != nil {
		logger.Fatalln(err)
	}
	return s
}

//store a Store implemented with blot as backend
type store struct {
	db     *bolt.DB
	logger *logrus.Logger
}

//PowerOn open a bolt instance
func (bs *store) init(conf *types.BoltConf, logger *logrus.Logger) error {

	bs.logger = logger

	var err error
	bs.db, err = bolt.Open(
		conf.Path,
		0600,
		&bolt.Options{
			Timeout: 1 * time.Second,
		})

	if err != nil {
		println(err.Error())
		return err
	}

	err = bs.db.Update(func(tx *bolt.Tx) error {

		_, e := tx.CreateBucketIfNotExists(BlobBucket)
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

//PowerOff close the blot instance
func (bs *store) PowerOff() error {
	return bs.db.Close()

}
