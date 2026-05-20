package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infrasphere/control-plane/backend/internal/config"
	"github.com/infrasphere/control-plane/backend/internal/database"
	"github.com/infrasphere/control-plane/backend/internal/providers/mock"
)

func TestHealth(t *testing.T) {
	srv := New(config.Load(), database.NewSeededStore(mock.New()))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

