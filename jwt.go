package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/pascaldekloe/jwt"
)

type JwtHelper struct {
	private_key *rsa.PrivateKey
	public_key  *rsa.PublicKey
}

type JwtSignatureTimeOut struct {
	error
}

func NewJwtHelper(public_key_path string, private_key_path string) (*JwtHelper, error) {
	public_key_file, err := os.OpenFile(public_key_path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	private_key_file, err := os.OpenFile(private_key_path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	public_key_byte, err := ioutil.ReadAll(public_key_file)
	if err != nil {
		return nil, err
	}
	private_key_byte, err := ioutil.ReadAll(private_key_file)
	if err != nil {
		return nil, err
	}

	public_key_block, _ := pem.Decode(public_key_byte)
	if public_key_block == nil {
		return nil, errors.New("can't decode public key")
	}
	private_key_block, _ := pem.Decode(private_key_byte)
	if private_key_block == nil {
		return nil, errors.New("can't decode private key")
	}

	public_key, err := x509.ParsePKIXPublicKey(public_key_block.Bytes)
	if err != nil {
		return nil, err
	}
	private_key, err := x509.ParsePKCS1PrivateKey(private_key_block.Bytes)
	if err != nil {
		return nil, err
	}

	return &JwtHelper{
		public_key:  public_key.(*rsa.PublicKey),
		private_key: private_key,
	}, nil
}

func (self JwtHelper) Verify(signature []byte) (*jwt.Claims, error) {
	claim, err := jwt.RSACheck(signature, self.public_key)
	if err != nil {
		return nil, err
	}
	valid := claim.Valid(time.Now())
	if !valid {
		return nil, &JwtSignatureTimeOut{}
	}
	return claim, nil
}

func (self JwtHelper) Sign(set map[string]interface{}, exp time.Time) ([]byte, error) {
	var claims jwt.Claims
	claims.Set = set
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(exp)
	return claims.RSASign(jwt.RS256, self.private_key, json.RawMessage(`{"type":"JWT"}`))
}
