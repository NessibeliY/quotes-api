package handlers

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	h := NewHandler(nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	buf := new(strings.Builder)
	_, err := io.Copy(buf, resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if buf.String() != "OK" {
		t.Errorf("expected response body OK, got %s", buf.String())
	}
}
