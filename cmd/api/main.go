package main

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/djalil83/A-Radius/internal/auth"
	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/djalil83/A-Radius/internal/customerportal"
	"github.com/djalil83/A-Radius/internal/dashboard/developer"
	"github.com/djalil83/A-Radius/internal/dashboard/pelanggan"
	"github.com/djalil83/A-Radius/internal/db/migrations"
	"github.com/djalil83/A-Radius/web"
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

	// RBAC runs behind the authenticated session.
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
		customerRoutes,
	)

	// Authentication wajib terjadi sebelum RBAC.
	customerAuthenticated := authHandler.RequireSession(
		protectedCustomerRoutes,
	)

	http.Handle(
		"/api/v1/customer/",
		customerAuthenticated,
	)

	// ------------------------------------------------------------
	// DEVELOPER SECURITY DASHBOARD
	// ------------------------------------------------------------
	// Authentication -> RBAC -> Developer Security Dashboard.
	//
	// Primary permission:
	//     security:scan
	//
	// The Developer dashboard is protected by the authenticated
	// session and server-side RBAC.

	developerHandler := developer.NewHandler()
	developerRoutes := developerHandler.Routes()

	protectedDeveloperRoutes := authzEngine.RequirePermissionHTTP(
		"security:scan",
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
					"developer dashboard authorization audit failed: %v",
					err,
				)
			}
		},
		developerRoutes,
	)

	// Authentication wajib terjadi sebelum RBAC.
	developerAuthenticated := authHandler.RequireSession(
		protectedDeveloperRoutes,
	)

	http.Handle(
		"/dashboard/developer/",
		developerAuthenticated,
	)

	// ------------------------------------------------------------
	// CUSTOMER DASHBOARD
	// ------------------------------------------------------------
	// Authentication -> RBAC -> role-specific customer dashboard.
	//
	// Permission:
	//     customer.portal.read
	//
	// The dashboard handler owns the embedded UI and the customer
	// dashboard endpoints. This keeps the customer dashboard behind
	// the same authentication/RBAC/audit boundary.
	pelangganHandler := pelanggan.NewHandler(customerHandler)

	pelangganRoutes := pelanggan.ProtectedRouter(
		pelangganHandler,
		authzEngine,
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
					"customer dashboard authorization audit failed: %v",
					err,
				)
			}
		},
	)

	// Authentication wajib terjadi sebelum RBAC.
	pelangganAuthenticated := authHandler.RequireSession(
		pelangganRoutes,
	)

	// Customer dashboard frontend asset.
	// dashboard.js hanya berisi kode frontend dan tidak mengekspos
	// data customer. API tetap dilindungi authentication + RBAC.
	http.Handle(
		"/dashboard/pelanggan/dashboard.js",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			data, err := fs.ReadFile(
				web.Assets,
				"dashboards/pelanggan/dashboard.js",
			)
			if err != nil {
				http.Error(
					w,
					"customer dashboard asset unavailable",
					http.StatusServiceUnavailable,
				)
				return
			}

			w.Header().Set(
				"Content-Type",
				"application/javascript; charset=utf-8",
			)
			w.Header().Set("Cache-Control", "no-store")

			_, _ = w.Write(data)
		}),
	)

	http.Handle(
		"/dashboard/pelanggan/",
		pelangganAuthenticated,
	)

	// Shared dashboard shell is non-sensitive frontend code.
	// It is served from the embedded web assets because dashboard.js
	// imports /shared/dashboard-shell.js.
	sharedAssets, err := fs.Sub(web.Assets, "shared")
	if err != nil {
		log.Fatalf("failed to initialize shared dashboard assets: %v", err)
	}

	http.Handle(
		"/shared/",
		http.StripPrefix(
			"/shared/",
			http.FileServer(http.FS(sharedAssets)),
		),
	)

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
