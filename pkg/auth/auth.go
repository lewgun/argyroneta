//package auth implement authorized.
package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const (
	goldenSalt = "something very salt"
)

//Authority implments authorized.
type Authority interface {

	//Verify verify the hash with the password & salt.
	Verify(privSalt, password, hash []byte) bool

	//HashAndSalt generate a new hash & salt pair from password.
	HashAndSalt(password []byte) ([]byte, []byte, error)
}

type auth struct{}

//Hash gen a hash from password &  private salt.
//alg: hash(pass + privSalt + goldenSalt)
func (auth) hash(privSalt, password []byte) ([]byte, error) {

	buf := &bytes.Buffer{}
	buf.Write(password)
	buf.Write(privSalt)
	buf.WriteString(goldenSalt)

	return bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost)

}

//Verify verify the hash with the password & salts.
func (auth) salt() ([]byte, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(p)))
	base64.StdEncoding.Encode(buf, p)

	return buf, nil
}

//Hash gen a hash & salt from password.
//alg: hash(pass + privSalt + goldenSalt)
func (au auth) HashAndSalt(password []byte) ([]byte, []byte, error) {

	salt, err := au.salt()
	if err != nil {
		return nil, nil, err
	}

	h, err := au.hash(salt, password)
	if err != nil {
		return nil, nil, err
	}
	return h, salt, nil

}

//Verify verify the hash with the password & salts.
func (auth) Verify(privSalt, password, hash []byte) bool {

	buf := &bytes.Buffer{}
	buf.Write(password)
	buf.Write(privSalt)
	buf.WriteString(goldenSalt)

	if err := bcrypt.CompareHashAndPassword(hash, buf.Bytes()); err != nil {
		return false
	}

	return true
}

//Auther is a helper for authorized.
var Auther Authority = auth{}
