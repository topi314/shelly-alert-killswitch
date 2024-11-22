package kswitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func New(cfg Config, version string) *Service {
	return &Service{
		cfg:     cfg,
		version: version,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type Service struct {
	cfg     Config
	version string
	server  *http.Server
	client  *http.Client
}

func (s *Service) Start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.Any("err", err))
	}
}

func (s *Service) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.Any("err", err))
	}
}

func (s *Service) SetupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", s.handleVersion)

	for _, cfg := range s.cfg.Configs {
		mux.HandleFunc(cfg.Endpoint, s.handleWebhook)
	}

	s.server = &http.Server{
		Addr:    s.cfg.Server.ListenAddr,
		Handler: mux,
	}
}

func (s *Service) handleWebhook(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "Received webhook", slog.String("path", r.URL.Path))

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	slog.Debug("Received event", slog.Any("event", event))

	for _, cfg := range s.cfg.Configs {
		if r.URL.Path == cfg.Endpoint {
			for _, alert := range event.Alerts {
				if alert.Labels.Contains(cfg.Labels) {
					s.executeKillSwitch(cfg, alert)
				}
			}
		}
	}
}

func (s *Service) executeKillSwitch(cfg KillSwitchConfig, alert Alert) {
	scheme := "https"
	if cfg.Insecure {
		scheme = "http"
	}

	var turn string
	switch alert.Status {
	case AlertStatusFiring:
		turn = "off"
	case AlertStatusResolved:
		turn = "on"
	default:
		slog.Error("Unknown alert status", slog.String("status", string(alert.Status)))
		return
	}

	rq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s/relay/%d?turn=%s", scheme, cfg.Address, cfg.Relay, turn), nil)
	if err != nil {
		slog.Error("Failed to create request", slog.Any("err", err))
		return
	}

	if cfg.Username != "" && cfg.Password != "" {
		rq.SetBasicAuth(cfg.Username, cfg.Password)
	}

	slog.Debug("Executing kill switch", slog.String("url", rq.URL.String()))

	rs, err := s.client.Do(rq)
	if err != nil {
		slog.Error("Failed to execute kill switch", slog.Any("err", err))
		return
	}
	defer func() {
		_ = rs.Body.Close()
	}()

	if rs.StatusCode != http.StatusOK {
		slog.Error("Error executing kill switch", slog.Any("status", rs.Status))
	}
}

func (s *Service) handleVersion(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(s.version)); err != nil {
		slog.Error("Failed to write version response", slog.Any("err", err))
	}
}
