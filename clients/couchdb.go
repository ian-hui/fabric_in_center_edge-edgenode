package clients

import (
	"fmt"
	"sync"

	"github.com/go-kivik/kivik/v4"
)

type couchdbClient struct {
	C  *kivik.Client
	Mu sync.Mutex
}

var (
	couchdbConns = &sync.Map{}
)

func InitCouchdb(kivik_addr string) error {
	client, err := kivik.New("couch", kivik_addr)
	if err != nil {
		return fmt.Errorf("init couchdb client error: %v", err)
	}
	kivikclient := new(couchdbClient)
	kivikclient.C = client
	couchdbConns.Store(kivik_addr, kivikclient)
	return nil
}

func GetCouchdb(kivik_addr string) (*couchdbClient, error) {
	pool, ok := couchdbConns.Load(kivik_addr)
	if !ok {
		InitCouchdb(kivik_addr)
		pool, _ = couchdbConns.Load(kivik_addr)
	}
	return pool.(*couchdbClient), nil
}
