package mux

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNativeContextMiddleware(t *testing.T) {
	withTimeout := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r := NewRouter()
	r.Handle("/path/{foo}", withTimeout(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := Vars(r)
		if vars["foo"] != "bar" {
			t.Fatal("Expected foo var to be set")
		}
	})))

	rec := NewRecorder()
	req := newRequest("GET", "/path/bar")
	r.ServeHTTP(rec, req)
}

func TestContextSetGet(t *testing.T) {
	tests := []struct {
		key   interface{}
		value interface{}
	}{
		{1, 2},
		{"test", 42},
	}

	req := newRequest("GET", "/path/bar")
	for _, test := range tests {
		contextSet(req, test.key, test.value)
	}
	for _, test := range tests {
		value := contextGet(req, test.key)
		if !reflect.DeepEqual(value, test.value) {
			t.Errorf(`Expected %+v got %+v`, test.value, value)
		}
	}
}
