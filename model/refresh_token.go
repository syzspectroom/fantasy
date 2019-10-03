package model

import (
	"math/rand"
	"time"
	"unsafe"
)

// RefreshToken db struct
type RefreshToken struct {
	ID             string `json:"_key,omitempty"`
	Token          string `json:"token"`
	UserAgent      string `json:"user_agent"`
	IP             string `json:"ip"`
	UserID         string `json:"user_id"`
	LastModifiedAt int64  `json:"last_modified_at"`
}

// FillWithMeta fill model with values from LoginMeta
func (rt *RefreshToken) FillWithMeta(loginMeta *LoginMeta) {
	rt.UserAgent = loginMeta.UserAgent
	rt.IP = loginMeta.IP
	rt.LastModifiedAt = time.Now().Unix()
	rt.UserID = loginMeta.UserID
	rt.GenerateRandomToken()
}

// GenerateRandomToken generates random refresh token
func (rt *RefreshToken) GenerateRandomToken() {
	// see https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const n = 40
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890./,][)(*&^%$#@_!+-"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	rt.Token = *(*string)(unsafe.Pointer(&b))
}
