package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/djalil83/A-Radius/internal/authn"
	"github.com/djalil83/A-Radius/internal/authz"
	admindashboard "github.com/djalil83/A-Radius/internal/dashboard/administrator"
	"github.com/djalil83/A-Radius/internal/subscriptionproduction"
	"github.com/djalil83/A-Radius/internal/subscriptionprofile"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := requiredEnv("DATABASE_URL")
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

	jwtConfig := authn.Config{
		Secret:   []byte(requiredEnv("JWT_SECRET")),
		Issuer:   requiredEnv("JWT_ISSUER"),
		Audience: requiredEnv("JWT_AUDIENCE"),
	}
	if err := jwtConfig.Validate(); err != nil {
		log.Fatalf("configure JWT middleware: %v", err)
	}

	authzEngine := authz.NewEngine(db)
	auditLogger := &authz.DBAuditLogger{DB: db}
	audit := func(ctx context.Context, principal *authz.Principal, permission string, allowed bool, status int, r *http.Request) {
		if err := auditLogger.AuthorizationDecision(ctx, principal, permission, allowed, status, r); err != nil {
			log.Printf("authorization audit: %v", err)
		}
	}

	profileProtected := subscriptionprofile.ProtectedRouter(
		subscriptionprofile.NewHandler(subscriptionprofile.NewRepository(db)),
		authzEngine,
		audit,
	)
	productionProtected := subscriptionproduction.ProtectedRouter(
		subscriptionproduction.NewHandler(),
		authzEngine,
		audit,
	)
	administratorProtected := admindashboard.ProtectedRouter(
		admindashboard.NewHandler(db),
		authzEngine,
		audit,
	)
	protected := chi.NewRouter()
	protected.Mount("/", profileProtected)
	protected.Mount("/", productionProtected)
	protected.Mount("/", administratorProtected)
	jwtHandler, err := jwtConfig.Middleware(protected)
	if err != nil {
		log.Fatalf("configure JWT middleware: %v", err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           jwtHandler,
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

func requiredEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}
