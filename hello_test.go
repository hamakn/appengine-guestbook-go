package guestbook

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
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

func TestSign(t *testing.T) {
	instance := aetestNewInstance(t)
	defer instance.Close()

	values := url.Values{}
	content := "aaaa42"
	values.Set("content", content)

	req, _ := instance.NewRequest(
		"POST",
		"/sign",
		strings.NewReader(values.Encode()),
	)
	// このrequestはformなのでcontent-typeを指定する
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := appengine.NewContext(req)

	res := httptest.NewRecorder()

	sign(res, req)

	if res.Code != http.StatusFound {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", res.Code)
	}

	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(ctx)).Limit(10)
	greetings := make([]Greeting, 0, 10)
	q.GetAll(ctx, &greetings)

	// greetingsのsizeが1であること
	if len := len(greetings); len != 1 {
		t.Fatalf("len(greetings) != 1, but %v", len)
	}

	// 最初のgreetingsのContentがpostしたものであること
	if g := greetings[0]; g.Content != content {
		t.Fatalf("greetings[0].Content != inputted_content, but %v", g.Content)
	}
}
