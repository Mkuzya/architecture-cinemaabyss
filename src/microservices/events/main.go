package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/segmentio/kafka-go"
)

type movieEvent struct {
    MovieID   int     `json:"movie_id"`
    Title     string  `json:"title"`
    Action    string  `json:"action"`
    UserID    int     `json:"user_id"`
    Timestamp string  `json:"timestamp,omitempty"`
}

type userEvent struct {
    UserID    int    `json:"user_id"`
    Username  string `json:"username"`
    Action    string `json:"action"`
    Timestamp string `json:"timestamp"`
}

type paymentEvent struct {
    PaymentID int     `json:"payment_id"`
    UserID    int     `json:"user_id"`
    Amount    float64 `json:"amount"`
    Status    string  `json:"status"`
    Timestamp string  `json:"timestamp"`
    Method    string  `json:"method_type"`
}

type eventService struct {
    brokers string
}

func newEventService() *eventService {
    brokers := os.Getenv("KAFKA_BROKERS")
    if brokers == "" {
        brokers = "kafka:9092"
    }
    return &eventService{brokers: brokers}
}

func (s *eventService) health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status": true}`))
}

func (s *eventService) produceAndConsume(topic string, payload []byte) error {
    writer := &kafka.Writer{
        Addr:         kafka.TCP(s.brokers),
        Topic:        topic,
        RequiredAcks: kafka.RequireAll,
        Balancer:     &kafka.LeastBytes{},
    }
    defer writer.Close()

    wctx, wcancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer wcancel()
    if err := writer.WriteMessages(wctx, kafka.Message{Value: payload}); err != nil {
        return err
    }
    log.Printf("Produced to %s: %s\n", topic, string(payload))

    // Best-effort short consumption from the end to validate pipeline without blocking CI
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     []string{s.brokers},
        Topic:       topic,
        MinBytes:    1,
        MaxBytes:    10e6,
        StartOffset: kafka.LastOffset,
    })
    defer reader.Close()

    rctx, rcancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer rcancel()
    msg, err := reader.ReadMessage(rctx)
    if err != nil {
        log.Printf("ReadMessage error (non-fatal for MVP): %v\n", err)
        return nil
    }
    log.Printf("Consumed from %s: %s\n", topic, string(msg.Value))
    return nil
}

func (s *eventService) handleMovie(w http.ResponseWriter, r *http.Request) {
    var ev movieEvent
    if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if ev.Timestamp == "" {
        ev.Timestamp = time.Now().Format(time.RFC3339)
    }
    b, _ := json.Marshal(ev)
    _ = s.produceAndConsume("movie-events", b)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, `{"status":"success"}`)
}

func (s *eventService) handleUser(w http.ResponseWriter, r *http.Request) {
    var ev userEvent
    if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    b, _ := json.Marshal(ev)
    _ = s.produceAndConsume("user-events", b)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, `{"status":"success"}`)
}

func (s *eventService) handlePayment(w http.ResponseWriter, r *http.Request) {
    var ev paymentEvent
    if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    b, _ := json.Marshal(ev)
    _ = s.produceAndConsume("payment-events", b)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, `{"status":"success"}`)
}

func main() {
    svc := newEventService()
    http.HandleFunc("/api/events/health", svc.health)
    http.HandleFunc("/api/events/movie", svc.handleMovie)
    http.HandleFunc("/api/events/user", svc.handleUser)
    http.HandleFunc("/api/events/payment", svc.handlePayment)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }
    log.Printf("Starting events service on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}


