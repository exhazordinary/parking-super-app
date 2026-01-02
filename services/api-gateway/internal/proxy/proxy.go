package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ServiceProxy handles request forwarding to backend services
type ServiceProxy struct {
	client *http.Client
}

func NewServiceProxy() *ServiceProxy {
	return &ServiceProxy{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// Forward proxies the request to the target service
func (p *ServiceProxy) Forward(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, `{"error":"invalid target URL"}`, http.StatusInternalServerError)
			return
		}

		// Build the full target URL
		proxyURL := *target
		proxyURL.Path = r.URL.Path
		proxyURL.RawQuery = r.URL.RawQuery

		// Create the proxy request
		proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, proxyURL.String(), r.Body)
		if err != nil {
			http.Error(w, `{"error":"failed to create request"}`, http.StatusInternalServerError)
			return
		}

		// Copy headers
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add forwarding headers
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		proxyReq.Header.Set("X-Forwarded-Host", r.Host)

		// Make the request
		resp, err := p.client.Do(proxyReq)
		if err != nil {
			log.Printf("proxy error: %v", err)
			http.Error(w, `{"error":"service unavailable"}`, http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Write response
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

// StripPrefix removes a prefix from the request path before forwarding
func (p *ServiceProxy) StripPrefix(prefix, targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		p.Forward(targetURL)(w, r)
	}
}
