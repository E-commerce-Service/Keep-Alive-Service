package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I am alive! Monitoring HTTP and DBs... ☕"))
	})

	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	targetsEnv := os.Getenv("TARGET_URLS")
	dbUrlsEnv := os.Getenv("DB_URLS")

	ticker := time.NewTicker(10 * time.Minute)

	runChecks(targetsEnv, dbUrlsEnv)

	for range ticker.C {
		runChecks(targetsEnv, dbUrlsEnv)
	}
}

func runChecks(httpTargets, dbTargets string) {
	log.Println("--- Starting Check Round ---")

	if httpTargets != "" {
		urls := strings.Split(httpTargets, ",")
		pingHttp(urls)
	} else {
		log.Println("⚠️ Nenhuma URL HTTP configurada.")
	}

	if dbTargets != "" {
		dbs := strings.Split(dbTargets, ",")
		pingDatabases(dbs)
	} else {
		log.Println("⚠️ Nenhum Banco de Dados configurado.")
	}
}

func pingHttp(targets []string) {
	client := http.Client{Timeout: 10 * time.Second}
	for _, url := range targets {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		resp, err := client.Get(url)
		if err != nil {
			log.Printf("❌ HTTP Falha: %s | Erro: %v", url, err)
			continue
		}
		resp.Body.Close()
		log.Printf("✅ HTTP Sucesso: %s | Status: %d", url, resp.StatusCode)
	}
}

func pingDatabases(connStrings []string) {
	for _, connStr := range connStrings {
		connStr = strings.TrimSpace(connStr)
		if connStr == "" {
			continue
		}

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("❌ DB Erro ao abrir driver: %v", err)
			continue
		}

		defer db.Close()

		err = db.Ping()
		if err != nil {
			safeLog := "Supabase/Postgres"
			if parts := strings.Split(connStr, "@"); len(parts) > 1 {
				safeLog = "..." + parts[1]
			}
			log.Printf("❌ DB Falha ao conectar em %s: %v", safeLog, err)
		} else {
			log.Printf("✅ DB Sucesso: Conexão ativa com o banco!")
		}
	}
}
