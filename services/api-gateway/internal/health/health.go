package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// ServiceHealth tracks health of backend services
type ServiceHealth struct {
	services map[string]string
	client   *http.Client
}

func NewServiceHealth(services map[string]string) *ServiceHealth {
	return &ServiceHealth{
		services: services,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type HealthStatus struct {
	Status   string                 `json:"status"`
	Services map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
}

// Handler returns the health check endpoint handler
func (h *ServiceHealth) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status := HealthStatus{
			Status:   "healthy",
			Services: make(map[string]ServiceInfo),
		}

		var wg sync.WaitGroup
		var mu sync.Mutex

		for name, url := range h.services {
			wg.Add(1)
			go func(name, url string) {
				defer wg.Done()

				info := h.checkService(ctx, url)

				mu.Lock()
				status.Services[name] = info
				if info.Status != "healthy" {
					status.Status = "degraded"
				}
				mu.Unlock()
			}(name, url)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func (h *ServiceHealth) checkService(ctx context.Context, baseURL string) ServiceInfo {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		return ServiceInfo{Status: "unhealthy"}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return ServiceInfo{Status: "unhealthy"}
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	if resp.StatusCode == http.StatusOK {
		return ServiceInfo{
			Status:  "healthy",
			Latency: latency.String(),
		}
	}

	return ServiceInfo{Status: "unhealthy"}
}
