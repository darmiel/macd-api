package keys

import (
	"context"
	"errors"
	"fmt"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type APIKey string

var (
	keys  []APIKey
	store limiter.Store
)

var (
	keyFile      = path.Join("data", "keys.txt")
	ctx          = context.Background()
	ErrNoFreeKey = errors.New("no free api key found")
)

func init() {
	var err error
	if store, err = memorystore.New(&memorystore.Config{
		Tokens:   4,
		Interval: time.Minute,
	}); err != nil {
		panic(err)
	}
	// load keys
	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	for _, l := range strings.Split(string(buf), "\n") {
		l = strings.TrimSpace(l)
		if len(l) != 16 {
			continue
		}
		keys = append(keys, APIKey(l))
	}
	fmt.Println("Loaded", len(keys), "api keys")
}

func (k APIKey) IsAvailable() bool {
	_, _, _, ok, err := store.Take(ctx, string(k))
	if err != nil {
		panic(err)
	}
	return ok
}

func (k APIKey) Invalidate() {
	// is there a better way?
	for {
		_, _, _, ok, err := store.Take(ctx, string(k))
		if err != nil {
			panic(err)
		}
		if !ok {
			break
		}
	}
}

func FindFreeKey() (APIKey, error) {
	for _, k := range keys {
		if k.IsAvailable() {
			return k, nil
		}
	}
	return "", ErrNoFreeKey
}
