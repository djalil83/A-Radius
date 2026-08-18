package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/djalil83/A-Radius/internal/auth"
	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/djalil83/A-Radius/internal/customerportal"
	"github.com/djalil83/A-Radius/internal/db/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error

	db, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Database connected successfully!")

	// ------------------------------------------------------------
	// DATABASE MIGRATIONS
	// ------------------------------------------------------------

	migrationDir := os.Getenv("MIGRATIONS_DIR")
	if migrationDir == "" {
		migrationDir = "database/postgresql/migrations"
	}

	migrationRunner := migrations.NewRunner(db, migrationDir)

	migrationCtx, migrationCancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer migrationCancel()

	if err := migrationRunner.Run(migrationCtx); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	log.Println("Database migrations applied successfully!")

	// ------------------------------------------------------------
	// AUTH
	// ------------------------------------------------------------

	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	authRoutes := http.NewServeMux()

	auth.RegisterRoutes(
		authRoutes,
		authHandler,
	)

	// /login and /logout do not require an existing session.
	// /me requires an authenticated session.
	authMe := authHandler.RequireSession(
		http.HandlerFunc(authHandler.Me),
	)

	authRoutes.Handle(
		"/api/v1/auth/me",
		authMe,
	)

	http.Handle(
		"/api/v1/auth/",
		authRoutes,
	)

	// ------------------------------------------------------------
	// AUTHORIZATION / RBAC
	// ------------------------------------------------------------

	authzEngine := authz.NewEngine(db)

	auditLogger := &authz.DBAuditLogger{
		DB: db,
	}

	// ------------------------------------------------------------
	// CUSTOMER PORTAL
	// ------------------------------------------------------------

	customerRepo := customerportal.NewRepository(db)
	customerService := customerportal.NewService(customerRepo)
	customerHandler := customerportal.NewHandler(customerService)

	customerRoutes := http.NewServeMux()

	customerportal.RegisterRoutes(
		customerRoutes,
		customerHandler,
	)

	// First authenticate the session and create the Principal.
	customerAuthenticated := authHandler.RequireSession(
		customerRoutes,
	)

	// Then enforce the RBAC permission.
	protectedCustomerRoutes := authzEngine.RequirePermissionHTTP(
		"customer.portal.read",
		func(
			ctx context.Context,
			principal *authz.Principal,
			permission string,
			allowed bool,
			status int,
			r *http.Request,
		) {
			if err := auditLogger.AuthorizationDecision(
				ctx,
				principal,
				permission,
				allowed,
				status,
				r,
			); err != nil {
				log.Printf(
					"authorization audit failed: %v",
					err,
				)
			}
		},
		customerAuthenticated,
	)

	http.Handle(
		"/api/v1/customer/",
		protectedCustomerRoutes,
	)

	// ------------------------------------------------------------
	// SYSTEM
	// ------------------------------------------------------------

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Server running on http://localhost:%s",
		port,
	)

	log.Printf(
		"Login API: http://localhost:%s/api/v1/auth/login",
		port,
	)

	log.Printf(
		"Auth API: http://localhost:%s/api/v1/auth/me",
		port,
	)

	log.Printf(
		"Customer API: http://localhost:%s/api/v1/customer/me",
		port,
	)

	log.Printf(
		"Dashboard API: http://localhost:%s/api/v1/customer/dashboard",
		port,
	)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("API server listening on :%s", port)
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignal)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("API server failed: %v", err)
		}

	case sig := <-shutdownSignal:
		log.Printf("Shutdown signal received: %v", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			_ = server.Close()
		}
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"text/plain; charset=utf-8",
	)

	_, _ = w.Write([]byte("pong"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		http.Error(
			w,
			"Database is not reachable",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/plain; charset=utf-8",
	)

	_, _ = w.Write([]byte("OK - Database connected"))
}
