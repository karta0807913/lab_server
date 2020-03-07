package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pascaldekloe/jwt"
)

func google_token(json_file string, scope string) (string, error) {
	file, err := os.OpenFile(json_file, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(file)
	var m map[string]interface{}
	decoder.Decode(&m)
	now := time.Now().Round(time.Second)
	to_time := now.Add(time.Second * 30)
	jwt_body := map[string]interface{}{
		"iss":   m["client_email"],
		"scope": scope,
		"aud":   "https://oauth2.googleapis.com/token",
		"exp":   int32(to_time.Unix()),
		"iat":   int32(now.Unix()),
	}

	var claims jwt.Claims
	claims.Set = jwt_body
	claims.Issued = jwt.NewNumericTime(now)
	claims.Expires = jwt.NewNumericTime(to_time)
	block, _ := pem.Decode([]byte(m["private_key"].(string)))
	if block == nil {
		return "", errors.New("can't decode private key")
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	token, err := claims.RSASign(jwt.RS256,
		private.(*rsa.PrivateKey), json.RawMessage(`{"type":"JWT"}`))
	if err != nil {
		return "", err
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token",
		url.Values{
			"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion":  {string(token)},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	decoder = json.NewDecoder(resp.Body)
	decoder.Decode(&m)
	if m["access_token"] == nil {
		return "", errors.New(m["error_description"].(string))
	}
	return m["access_token"].(string), nil
}
