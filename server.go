package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// item é o recurso servido por GET /items/{id}.
// As struct tags definem o nome de cada campo no JSON.
type item struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"` // dinheiro em centavos: float tem erro de arredondamento
}

type healthResponse struct {
	Status        string  `json:"status"`
	UptimeSeconds float64 `json:"uptime_seconds"`
}

type metricsResponse struct {
	RequestsTotal int64 `json:"requests_total"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// server guarda o estado compartilhado pelos handlers.
type server struct {
	startTime time.Time
	// atomic.Int64 porque cada requisição roda na sua própria goroutine:
	// um int comum com ++ teria condição de corrida (race condition).
	requests atomic.Int64
	// items é só lido depois de criado, então não precisa de mutex.
	items map[int]item
}

// newServer recebe a hora de início como parâmetro para os testes
// poderem simular um servidor que já está no ar há algum tempo.
func newServer(startTime time.Time) *server {
	return &server{
		startTime: startTime,
		items: map[int]item{
			1: {ID: 1, Name: "Teclado mecânico", PriceCents: 29990},
			2: {ID: 2, Name: "Mouse sem fio", PriceCents: 8990},
			3: {ID: 3, Name: "Monitor 24 polegadas", PriceCents: 89900},
		},
	}
}

// routes liga cada rota ao seu handler. Desde o Go 1.22 o padrão aceita
// método ("GET ...") e curingas ("{id}"); outro método recebe 405.
func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /metrics", s.handleMetrics)
	mux.HandleFunc("GET /items/{id}", s.handleGetItem)
	return s.countRequests(mux)
}

// countRequests é um middleware: embrulha outro handler e roda código
// antes dele. Aqui, conta toda requisição (inclusive a do próprio /metrics).
func (s *server) countRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:        "ok",
		UptimeSeconds: time.Since(s.startTime).Seconds(),
	})
}

func (s *server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, metricsResponse{RequestsTotal: s.requests.Load()})
}

func (s *server) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "id deve ser um inteiro positivo"})
		return
	}

	it, ok := s.items[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "item não encontrado"})
		return
	}

	writeJSON(w, http.StatusOK, it)
}

// writeJSON centraliza o formato das respostas. O cabeçalho precisa ser
// definido antes de WriteHeader, que envia status e cabeçalhos ao cliente.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("erro ao escrever JSON: %v", err)
	}
}
