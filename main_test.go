package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz: got %d, want 200", rec.Code)
	}
}

func TestVersion(t *testing.T) {
	version = "1.2.3"
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	var v versionInfo
	if err := json.NewDecoder(rec.Body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v.Version != "1.2.3" {
		t.Fatalf("version: got %q, want 1.2.3", v.Version)
	}
}

func TestRootHello(t *testing.T) {
	t.Setenv("LADDER_FAIL_RATE", "0")
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("root: got %d, want 200", rec.Code)
	}
}

func TestFailRateParsing(t *testing.T) {
	cases := map[string]float64{"": 0, "abc": 0, "1.5": 0, "0.25": 0.25}
	for in, want := range cases {
		t.Setenv("LADDER_FAIL_RATE", in)
		if got := failRate(); got != want {
			t.Errorf("failRate(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestEnvOr(t *testing.T) {
	t.Setenv("PORT", "")
	if envOr("PORT", "8080") != "8080" {
		t.Fatal("default not applied")
	}
}
