/******************************************************************************
* FILENAME:      utils.go
*
* AUTHORS:       Qiu Siyu START DATE: 周三 7月 13 2022
*
* LAST MODIFIED: 星期三, 七月 13th 2022, 上午10：40
*
* CONTACT:       siyu.qiu@smartmore.com
******************************************************************************/

package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	runes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken = errors.Unauthorized(reason, "JWT token is missing")
)

func MD5Hash(strs ...string) string {
	res := ""
	for _, str := range strs {
		res = fmt.Sprint(res, str)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(res)))
}

func RandomCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

// RandomString generates a random string with length @size
// DefaultSaltLength = 16
func RandomString(size int) string {
	if size == 0 {
		size = 16
	}
	b := make([]byte, size)
	for i := range b {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		b[i] = runes[r.Int63()%int64(len(runes))]
	}
	return string(b)
}

func GetTokenFromContext(ctx context.Context) (string, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return "", fmt.Errorf("transport empty")
	}
	tk := ts.RequestHeader().Get("Authorization")
	auths := strings.SplitN(tk, " ", 2)
	if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
		return "", ErrMissingJwtToken
	}
	jwtToken := auths[1]
	return jwtToken, nil
}

func SliceContainsAny[T string](arr []T, target T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
