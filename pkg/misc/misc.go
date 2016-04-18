package misc

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Unused(...interface{}) {

}

//gzip
func Zip(in []byte) ([]byte, error) {

	var (
		err error
	)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(in)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//gunzip
func Unzip(in []byte) ([]byte, error) {

	rdr := bytes.NewReader(in)
	r, err := gzip.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	r.Close()
	return out, nil
}
