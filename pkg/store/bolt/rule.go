package bolt

import (
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/rule"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	"github.com/pguerna/ffjson/ffjson"
)

func (bs *BoltStore) AddRule(r *rule.Rule) error {
	return bs.UpdateRule(r)
}

func (bs *BoltStore) UpdateRule(r *rule.Rule) error {

	data, err := ffjson.Marshal(r)
	if err != nil {
		return err
	}
	err = bs.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket(RuleBucket)
		return bu.Put([]byte(r.Ident), data)
	})
	return err

}

func (bs *BoltStore) DeleteRule(key string) error {

	// Delete the key in a different write transaction.
	err := bs.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(RuleBucket).Delete([]byte(key))
	})
	return err

}

func (bs *BoltStore) Rule(key string) (*rule.Rule, error) {

	var (
		data []byte
		r    *rule.Rule = &rule.Rule{}
	)

	err := bs.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket(RuleBucket)
		data = bu.Get([]byte(key))
		if len(data) == 0 {
			//color.Red(err.Error())
			return errutil.ErrNotFound
		}
		return nil
	})
	if err != nil {
		return nil, false
	}

	err = ffjson.Unmarshal(data, r)
	if err != nil {
		color.Red(err.Error())
		return nil, false
	}

	return r, nil

}

func (bs *BoltStore) Rules() (map[string]*rule.Rule, error) {
	m := map[string]*rule.Rule{}
	var err error

	bs.db.View(func(tx *bolt.Tx) error {
		bu := tx.Bucket(RuleBucket)
		err = bu.ForEach(func(k, v []byte) error {
			var ru = &rule.Rule{}
			e := ffjson.Unmarshal(v, ru)
			if e != nil {
				color.Red(err.Error())
				return e
			}
			m[ru] = string(k)
			return nil
		})
		return err
	})
	return m, err

}
