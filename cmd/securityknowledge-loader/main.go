package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/djalil83/A-Radius/internal/securityknowledge"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	ctx := context.Background()

	repository := securityknowledge.NewRepository(db)
	loader := securityknowledge.NewLoader(repository)

	knowledge, err := loader.LoadManifestFile(
		ctx,
		"security/knowledge/v1/manifest.json",
		"security/knowledge/v1/manifest.json",
	)
	if err != nil {
		log.Fatalf("load knowledge: %v", err)
	}

	fmt.Printf(
		"draft created: key=%s version=%s status=%s hash=%s\n",
		knowledge.KnowledgeKey,
		knowledge.Version,
		knowledge.Status,
		knowledge.ContentHash,
	)
}
