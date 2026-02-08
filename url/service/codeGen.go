package service

import (
	"crypto/rand"
	"errors"
	"io"
)

const (
	codeLen    = 7
	maxRetries = 5
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 62
)

var ErrCollision = errors.New("short code collision after max retries")

// ExistsFn должна вернуть true, если код уже занят (например, есть в БД).
type ExistsFn func(code string) (bool, error)

// GenerateShortCode генерирует 7-символьный код из [a-zA-Z0-9] с crypto/rand,
// и проверяет коллизию через exists. Делает до 5 попыток.
func GenerateShortCode(exists ExistsFn) (string, error) {
	if exists == nil {
		return "", errors.New("exists func is nil")
	}

	for i := 0; i < maxRetries; i++ {
		code, err := randomBase62String(rand.Reader, codeLen)
		if err != nil {
			return "", err
		}

		taken, err := exists(code)
		if err != nil {
			return "", err
		}
		if !taken {
			return code, nil
		}
	}

	return "", ErrCollision
}

// randomBase62String равномерно выбирает символы из charset через rejection sampling.
func randomBase62String(r io.Reader, n int) (string, error) {
	if n <= 0 {
		return "", errors.New("invalid length")
	}

	out := make([]byte, n)
	i := 0

	const max = 62 * 4 // 248

	var buf [1]byte
	for i < n {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return "", err
		}
		b := buf[0]
		if b >= max {
			continue
		}
		out[i] = charset[int(b)%len(charset)]
		i++
	}

	return string(out), nil
}
