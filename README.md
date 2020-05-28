# DLock [![GoDoc](https://godoc.org/github.com/4vn/dlock?status.svg)](https://pkg.go.dev/github.com/4vn/dlock)

Distributed lock using redis

### Import

``` go
import "github.com/4vn/dlock"
```

### Usage

``` go
l := dlock.New("127.0.0.1:6379")
if lockId, err := l.Lock("user1", 100*time.Millisecond); err != nil {
	log.Println("Lock failed")
}
defer l.Unlock("user1", lockId)

// Do something
```

### License

MIT
