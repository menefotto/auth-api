package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Thing struct{}

func (t Thing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok")
}

func TestToJson(t *testing.T) {
	tests := []struct {
		Url     string
		Payload []byte
		Type    string
	}{
		{
			Url:     "/",
			Payload: []byte(`{"Username":"wind85"}`),
			Type:    "application/json",
		},
		{
			Url:     "/",
			Payload: []byte(`{"Nope"}`),
			Type:    "application/json",
		},
	}

	test := &Thing{}
	s := httptest.NewServer(ToJson(test))
	defer s.Close()

	for _, ts := range tests {
		content, err := json.Marshal(ts.Payload)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.Post(s.URL, ts.Type, bytes.NewBuffer(content))
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Fatal("Request gone wrong, status: ", resp.StatusCode)
		}
	}
}
