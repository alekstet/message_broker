package models

import (
	"sync"
)

type Data struct {
	Sync_ chan struct{}
	Cnt   int
	Url   string
	Queue []string
	Wg    sync.WaitGroup
	Mu    sync.Mutex
}
