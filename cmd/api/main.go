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

	bootstrap "github.com/FabioRocha231/saas-core/internal/infra/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// ===== Config básica =====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// gin.SetMode(gin.ReleaseMode)

	// ===== Router =====
	r := gin.New()
	
	// Rotas
	bootstrap.RegisterRoutes(r)

	// ===== HTTP Server com timeouts =====
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// ===== Start server em goroutine =====
	go func() {
		log.Printf("[saas-core] listening on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[saas-core] listen error: %v", err)
		}
	}()

	// ===== Graceful shutdown =====
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("[saas-core] shutdown signal received: %s\n", sig.String())

	// Tempo máximo para encerrar requests em andamento
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[saas-core] shutdown error: %v\n", err)
		// Se der erro, força fechamento
		_ = srv.Close()
	}

	log.Println("[saas-core] server stopped gracefully")
}