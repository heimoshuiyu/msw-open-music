package api

import "testing"

func TestAPI(t *testing.T) {
	api, err := NewAPI("/tmp/test.sqlite3", ":8080")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Fatal(api.Server.ListenAndServe())
}
