package guestbook

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine/aetest"
)

func aetestNewInstance(t *testing.T) aetest.Instance {
	opt := aetest.Options{StronglyConsistentDatastore: true}
	instance, err := aetest.NewInstance(&opt)
	if err != nil {
		t.Fatalf("Failed to create aetest instance: %v", err)
	}
	return instance
}

func TestRoot(t *testing.T) {
	instance := aetestNewInstance(t)
	defer instance.Close()

	req, _ := instance.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	root(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", res.Code)
	}
}
