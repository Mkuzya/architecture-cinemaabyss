package main

import (
    "io"
    "log"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"
)

type upstreamSelector struct {
    monolithURL      string
    moviesServiceURL string
    gradualMigration bool
    moviesPercent    int
    httpClient       *http.Client
}

func newUpstreamSelector() *upstreamSelector {
    gradual := strings.EqualFold(os.Getenv("GRADUAL_MIGRATION"), "true")
    percentStr := os.Getenv("MOVIES_MIGRATION_PERCENT")
    percent := 0
    if p, err := strconv.Atoi(percentStr); err == nil {
        if p < 0 {
            p = 0
        }
        if p > 100 {
            p = 100
        }
        percent = p
    }

    client := &http.Client{Timeout: 30 * time.Second}

    return &upstreamSelector{
        monolithURL:      getenvDefault("MONOLITH_URL", "http://monolith:8080"),
        moviesServiceURL: getenvDefault("MOVIES_SERVICE_URL", "http://movies-service:8081"),
        gradualMigration: gradual,
        moviesPercent:    percent,
        httpClient:       client,
    }
}

func getenvDefault(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func (u *upstreamSelector) chooseMoviesUpstream() string {
    if !u.gradualMigration {
        return u.monolithURL
    }
    n := rand.Intn(100)
    if n < u.moviesPercent {
        return u.moviesServiceURL
    }
    return u.monolithURL
}

func (u *upstreamSelector) proxy(w http.ResponseWriter, r *http.Request, upstreamBase string) {
    targetURL, err := url.Parse(upstreamBase)
    if err != nil {
        http.Error(w, "invalid upstream", http.StatusBadGateway)
        return
    }

    // Preserve path and query
    targetURL.Path = singleJoinPath(targetURL.Path, r.URL.Path)
    targetURL.RawQuery = r.URL.RawQuery

    req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String(), r.Body)
    if err != nil {
        http.Error(w, "failed to create upstream request", http.StatusBadGateway)
        return
    }
    // copy headers
    req.Header = r.Header.Clone()

    resp, err := u.httpClient.Do(req)
    if err != nil {
        http.Error(w, "upstream request failed", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // copy response headers and status
    for k, vv := range resp.Header {
        for _, v := range vv {
            w.Header().Add(k, v)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func singleJoinPath(a, b string) string {
    as := strings.TrimRight(a, "/")
    bs := strings.TrimLeft(b, "/")
    return as + "/" + bs
}

func main() {
    rand.Seed(time.Now().UnixNano())
    selector := newUpstreamSelector()

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Movies go through Strangler Fig
    http.HandleFunc("/api/movies", func(w http.ResponseWriter, r *http.Request) {
        upstream := selector.chooseMoviesUpstream()
        selector.proxy(w, r, upstream)
    })
    http.HandleFunc("/api/movies/", func(w http.ResponseWriter, r *http.Request) {
        upstream := selector.chooseMoviesUpstream()
        selector.proxy(w, r, upstream)
    })

    // All other APIs go to monolith by default
    http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
        selector.proxy(w, r, selector.monolithURL)
    })

    port := getenvDefault("PORT", "8000")
    log.Printf("Starting proxy on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
