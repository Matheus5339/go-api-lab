// Comando go-api-lab: uma API HTTP pequena, só com a biblioteca padrão do Go.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// O Cloud Run (e muitas plataformas) informa a porta pela variável PORT.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: newServer(time.Now()).routes(),
		// Sem timeout, um cliente lento pode segurar conexões para sempre
		// (ataque "Slowloris").
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ctx é cancelado quando o processo recebe Ctrl+C ou SIGTERM
	// (o Cloud Run manda SIGTERM antes de desligar o contêiner).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// O servidor roda numa goroutine para o main poder esperar o sinal.
	go func() {
		log.Printf("ouvindo em %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	<-ctx.Done()

	// Graceful shutdown: para de aceitar conexões novas e espera
	// as requisições em andamento terminarem (até 10s).
	log.Println("desligando...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro ao desligar: %v", err)
	}
}
