package server

import "testing"

func TestGetAll(t *testing.T) {
	want := "a"
	got := "a"
	if got != want {
		t.Errorf("Got %s, Want %s", got, want)
	}
}
