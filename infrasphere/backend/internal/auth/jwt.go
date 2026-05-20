package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims struct {
	Sub   string   `json:"sub"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	Exp   int64    `json:"exp"`
}

func Sign(secret string, claims Claims) (string, error) {
	claims.Exp = time.Now().Add(12 * time.Hour).Unix()
	header := enc(map[string]string{"alg": "HS256", "typ": "JWT"})
	body := enc(claims)
	msg := header + "." + body
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(msg))
	return msg + "." + base64.RawURLEncoding.EncodeToString(sig.Sum(nil)), nil
}

func Verify(secret, token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	msg := parts[0] + "." + parts[1]
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(msg))
	expected := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("invalid signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return nil, err
	}
	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func enc(v interface{}) string {
	b, _ := json.Marshal(v)
	return base64.RawURLEncoding.EncodeToString(b)
}

