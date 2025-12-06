package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I am alive and waking up others! ☕"))
	})

	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	targetsEnv := os.Getenv("TARGET_URLS")
	if targetsEnv == "" {
		log.Println("⚠️ Nenhuma URL configurada em TARGET_URLS")
	} else {
		targets := strings.Split(targetsEnv, ",")
		ticker := time.NewTicker(10 * time.Minute)
		pingHosts(targets)

		for range ticker.C {
			pingHosts(targets)
		}
	}

	select {}
}

func pingHosts(targets []string) {
	log.Println("--- Starting Ping Round ---")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	for _, url := range targets {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		resp, err := client.Get(url)
		if err != nil {
			log.Printf("❌ Falha ao acordar %s: %v", url, err)
			continue
		}
		resp.Body.Close()
		log.Printf("✅ Sucesso %s: Status %d", url, resp.StatusCode)
	}
}
