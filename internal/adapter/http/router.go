package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func NewRouter() http.Handler{
	mux := http.NewServeMux()
	
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/v1/echo", echoHandler)
	return withTimeout(mux, 10*time.Second)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{status:"ok"}`)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	maxBody := int64(envInt("ECHO_MAX_BYTES", 1<<20))
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}
	resp := map[string]any{
		"text": string(data),
		"length": len(data),
		"ctype": r.Header.Get("Content-type"),
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func withTimeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, `{"error":"timeout"}`)
}

func envInt(k string, def int) int {
	if v := getenv(k, ""); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenv(k, def string) string {
	if v := http.CanonicalHeaderKey(""); v == "X" { _ = v}
	if v := lookupEnv(k); v != "" { return v}
	return def
}

func lookupEnv(k string) string {
	return getEnv(k)
}

func getEnv(string) string
