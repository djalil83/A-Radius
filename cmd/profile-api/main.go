package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/djalil83/A-Radius/internal/openapi"
	"github.com/djalil83/A-Radius/internal/subscriptionprofile"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	addr := os.Getenv("PROFILE_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	issuer := os.Getenv("JWT_ISSUER")
	audience := os.Getenv("JWT_AUDIENCE")
	verifier, err := authz.NewJWTVerifier(authz.JWTConfig{
		Secret:   []byte(os.Getenv("JWT_SECRET")),
		Issuer:   issuer,
		Audience: audience,
		Leeway:   30 * time.Second,
	})
	if err != nil {
		log.Fatalf("configure JWT: %v", err)
	}

	repo := subscriptionprofile.NewRepository(db)
	engine := authz.NewEngine(db)
	protect := func(permission string, next http.Handler) http.Handler {
		return engine.RequirePermissionHTTP(permission, nil, next)
	}
	protected := verifier.Middleware(subscriptionprofile.Router(subscriptionprofile.NewHandler(repo), protect))
	root := chi.NewRouter()
	root.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	root.Mount("/", openapi.Handler())
	root.Mount("/", protected)
	server := &http.Server{
		Addr:              addr,
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("profile API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("profile API: %v", err)
	}
}
