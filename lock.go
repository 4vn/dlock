package dlock

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	lockScript   = redis.NewScript(`return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])`)
	unlockScript = redis.NewScript(`if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end`)
)

type Lock struct {
	c *redis.Client
}

// New creates Lock object
func New(addr string) *Lock {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: "", DB: 0})
	return &Lock{c: client}
}

// Lock attempts to put a lock on the key for a specified duration.
// returns error if failed.
func (l *Lock) Lock(key string, timeout time.Duration) (string, error) {
	id := randStr(10)
	res, err := lockScript.Run(context.Background(), l.c, []string{key}, id, timeout.Milliseconds()).Text()
	if err != nil && err != redis.Nil {
		return "", err
	}

	if res != "OK" {
		return "", errors.New("Lock failed")
	}

	return id, nil
}

// Unlock attempts to remove the lock on a key so long as the value matches.
// returns error if failed.
func (l *Lock) Unlock(key, id string) error {
	res, err := unlockScript.Run(context.Background(), l.c, []string{key}, id).Int()
	if err != redis.Nil && err != nil {
		return err
	}

	if res != 1 {
		return errors.New("Unlock failed")
	}

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
