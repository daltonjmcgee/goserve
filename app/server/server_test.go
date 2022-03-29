package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleStaticDirectory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/static/", nil)
	w := httptest.NewRecorder()
	handleStatic(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// Missing config info
	if string(data) != "Some dipshit deleted the default 404 and didn't replace it. At any rate, your page wasn't found." {
		t.Errorf("expected ABC got %v", string(data))
	}
}

func TestHandleStaticFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/static/main.js", nil)
	w := httptest.NewRecorder()
	handleStatic(w, req)
	res := w.Result()
	fmt.Println(res)
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "" {
		t.Errorf("expected Nothing got %v", string(data))
	}
}
