package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alekstet/message_broker/endpoint"
)

func TestQueue(t *testing.T) {
	d := endpoint.New("/color")
	resp, _ := json.Marshal("red")
	resp1, _ := json.Marshal("red1")

	handler := http.HandlerFunc(d.Endpoint)

	go t.Run("GET without timeout", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/color", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, "404 Not Found", w.Result().Status)
	})

	go t.Run("GET without PUT", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/color?timeout=2", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, "404 Not Found", w.Result().Status)
	})

	time.Sleep(time.Second * 3)

	go t.Run("GET with PUT 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/color?timeout=10", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(resp), w.Body.Bytes())
	})

	time.Sleep(time.Second * 1)

	go t.Run("GET with PUT 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/color?timeout=15", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(resp1), w.Body.Bytes())
	})

	go t.Run("PUT 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/color?v=red", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(nil), w.Body.Bytes())
	})

	time.Sleep(time.Second * 1)

	go t.Run("PUT 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/color?v=red1", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(nil), w.Body.Bytes())
	})
	time.Sleep(time.Second * 3)
}
