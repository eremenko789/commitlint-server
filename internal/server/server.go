package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// Server represents the webhook server
type Server struct {
	config *lint.Config
	linter *lint.Linter
}

// Run starts the webhook server
func Run() error {
	// Load configuration
	conf, err := config.LookupAndParse()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Validate server configuration
	if err := validateServerConfig(&conf.Server); err != nil {
		return fmt.Errorf("ошибка конфигурации сервера: %w", err)
	}

	// Create linter
	linter, err := config.NewLinter(conf)
	if err != nil {
		return fmt.Errorf("ошибка создания линтера: %w", err)
	}

	server := &Server{
		config: conf,
		linter: linter,
	}

	return server.start()
}

// validateServerConfig validates server configuration
func validateServerConfig(cfg *lint.ServerConfig) error {
	if cfg.Address == "" {
		cfg.Address = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.GiteaURL == "" {
		return fmt.Errorf("gitea_url обязателен")
	}
	if cfg.GiteaToken == "" {
		return fmt.Errorf("gitea_token обязателен")
	}
	return nil
}

// start starts the HTTP server
func (s *Server) start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhook)
	mux.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)
	
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Получен сигнал остановки, завершаем работу...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Ошибка при остановке сервера: %v", err)
		}
	}()

	log.Printf("Сервер запущен на %s", addr)
	
	if s.config.Server.CertFile != "" && s.config.Server.KeyFile != "" {
		return server.ListenAndServeTLS(s.config.Server.CertFile, s.config.Server.KeyFile)
	}
	
	return server.ListenAndServe()
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleWebhook handles Gitea webhook requests
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела запроса: %v", err)
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	// Verify webhook signature if secret is configured
	if s.config.Server.WebhookSecret != "" {
		if !s.verifySignature(r, body) {
			log.Println("Неверная подпись webhook")
			http.Error(w, "Неверная подпись", http.StatusUnauthorized)
			return
		}
	}

	eventType := r.Header.Get("X-Gitea-Event")
	if s.config.Server.Debug {
		log.Printf("Получен webhook: тип=%s", eventType)
	}

	switch eventType {
	case "pull_request":
		s.handlePullRequestWebhook(w, body)
	default:
		if s.config.Server.Debug {
			log.Printf("Игнорируем событие типа: %s", eventType)
		}
		w.WriteHeader(http.StatusOK)
	}
}

// verifySignature verifies the webhook signature
func (s *Server) verifySignature(r *http.Request, body []byte) bool {
	signature := r.Header.Get("X-Gitea-Signature")
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.config.Server.WebhookSecret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// handlePullRequestWebhook handles pull request webhooks
func (s *Server) handlePullRequestWebhook(w http.ResponseWriter, body []byte) {
	var webhook PullRequestWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Printf("Ошибка парсинга webhook pull request: %v", err)
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Only process opened and synchronized events
	if webhook.Action != "opened" && webhook.Action != "synchronize" {
		if s.config.Server.Debug {
			log.Printf("Игнорируем действие pull request: %s", webhook.Action)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if s.config.Server.Debug {
		log.Printf("Обрабатываем pull request #%d: %s", 
			webhook.PullRequest.Number, webhook.PullRequest.Title)
	}

	// Process the pull request
	go s.processPullRequest(&webhook)
	
	w.WriteHeader(http.StatusOK)
}