package clients

import (
	"sync"

	"github.com/go-kivik/kivik/v4"
)

type couchdbClient struct {
	C  *kivik.Client
	Mu sync.Mutex
}

var (
	mu           sync.Mutex
	couchdbConns *couchdbClient
)

// func InitCouchdb(kivik_addr string) error {
// 	client, err := kivik.New("couch", kivik_addr)
// 	if err != nil {
// 		return fmt.Errorf("init couchdb client error: %v", err)
// 	}
// 	kivikclient := new(couchdbClient)
// 	kivikclient.C = client
// 	couchdbConns.Store(kivik_addr, kivikclient)
// 	return nil
// }

func GetCouchdb(kivik_addr string) (*couchdbClient, error) {
	if couchdbConns == nil {
		mu.Lock()
		defer mu.Unlock()
		if couchdbConns == nil {
			client, err := kivik.New("couch", kivik_addr)
			if err != nil {
				return nil, err
			}
			kivikclient := new(couchdbClient)
			kivikclient.C = client
			kivikclient.Mu = sync.Mutex{}
			couchdbConns = kivikclient
		}
	}
	return couchdbConns, nil
}
