package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// do envia uma requisição direto para o handler, sem abrir porta de rede:
// httptest.NewRecorder faz o papel do http.ResponseWriter e guarda a resposta.
func do(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// decode converte o corpo JSON da resposta para o tipo T.
func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(rec.Body).Decode(&v); err != nil {
		t.Fatalf("corpo não é JSON válido: %v (corpo: %q)", err, rec.Body.String())
	}
	return v
}

func TestHealth(t *testing.T) {
	// Servidor "no ar" há 10 segundos.
	h := newServer(time.Now().Add(-10 * time.Second)).routes()

	rec := do(t, h, http.MethodGet, "/health")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quero %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, quero application/json", ct)
	}
	body := decode[healthResponse](t, rec)
	if body.Status != "ok" {
		t.Errorf("status = %q, quero ok", body.Status)
	}
	if body.UptimeSeconds < 10 {
		t.Errorf("uptime_seconds = %v, quero >= 10", body.UptimeSeconds)
	}
}

func TestMetricsCountsRequests(t *testing.T) {
	h := newServer(time.Now()).routes()

	do(t, h, http.MethodGet, "/health")
	do(t, h, http.MethodGet, "/items/1")
	do(t, h, http.MethodGet, "/items/abc")
	rec := do(t, h, http.MethodGet, "/metrics")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quero %d", rec.Code, http.StatusOK)
	}
	// 3 requisições anteriores + a própria chamada a /metrics.
	if got := decode[metricsResponse](t, rec).RequestsTotal; got != 4 {
		t.Errorf("requests_total = %d, quero 4", got)
	}
}

// Teste em tabela (table-driven test): cada linha é um caso, e t.Run
// cria um subteste com nome, que aparece separado na saída do go test -v.
func TestGetItem(t *testing.T) {
	h := newServer(time.Now()).routes()

	tests := []struct {
		name     string
		path     string
		wantCode int
		wantName string // só conferido quando wantCode == 200
	}{
		{"item existente", "/items/1", http.StatusOK, "Teclado mecânico"},
		{"id não numérico", "/items/abc", http.StatusBadRequest, ""},
		{"id zero", "/items/0", http.StatusBadRequest, ""},
		{"id negativo", "/items/-5", http.StatusBadRequest, ""},
		{"id inexistente", "/items/999", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, h, http.MethodGet, tt.path)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, quero %d (corpo: %s)", rec.Code, tt.wantCode, rec.Body.String())
			}
			if tt.wantCode == http.StatusOK {
				if got := decode[item](t, rec).Name; got != tt.wantName {
					t.Errorf("name = %q, quero %q", got, tt.wantName)
				}
				return
			}
			if decode[errorResponse](t, rec).Error == "" {
				t.Error("resposta de erro sem o campo error")
			}
		})
	}
}

func TestWrongMethodReturns405(t *testing.T) {
	h := newServer(time.Now()).routes()

	rec := do(t, h, http.MethodPost, "/health")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, quero %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
