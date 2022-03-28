package endpoint

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/alekstet/message_broker/models"
)

type D models.Data

func New(url string) *D {
	cnt := 0
	queue := []string{}
	sync_ := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	return &D{
		Sync_: sync_,
		Cnt:   cnt,
		Url:   url,
		Queue: queue,
		Wg:    wg,
		Mu:    mu,
	}
}

func (d *D) Timeout(c chan struct{}, timeout int, w http.ResponseWriter) {
	go func() {
		for {
			if len(d.Queue) > 0 {
				c <- struct{}{}
				return
			}
		}
	}()
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		w.WriteHeader(404)
		d.Wg.Done()
	case <-c:
		jsonResp, err := json.Marshal(d.Queue[0])
		if err != nil {
			w.WriteHeader(404)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
		d.Queue = d.Queue[1:]
		d.Sync_ <- struct{}{}
		d.Wg.Done()
	}
}

func (d *D) Endpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		timeout := r.URL.Query().Get("timeout")
		if timeout != "" {
			d.Cnt++
			d.Wg.Add(1)
			timeout_int, _ := strconv.Atoi(timeout)
			ch := make(chan struct{})
			if d.Cnt == 1 {
				go d.Timeout(ch, timeout_int, w)
			}
			if d.Cnt > 1 {
				<-d.Sync_
				go d.Timeout(ch, timeout_int, w)
			}
			d.Wg.Wait()
		} else {
			if len(d.Queue) > 0 {
				jsonResp, err := json.Marshal(d.Queue[0])
				if err != nil {
					w.WriteHeader(404)
				}
				d.Mu.Lock()
				defer d.Mu.Unlock()
				d.Queue = d.Queue[1:]
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResp)
			} else {
				w.WriteHeader(404)
			}
		}
	case "PUT":
		if r.URL.Query().Get("v") == "" {
			w.WriteHeader(400)
		} else {
			d.Mu.Lock()
			defer d.Mu.Unlock()
			d.Queue = append(d.Queue, r.URL.Query().Get("v"))
		}
	}
}
