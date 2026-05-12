package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
    redisClient *redis.Client
    mongoClient *mongo.Client
    logger      *log.Logger
)

type HealthResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

func main() {
    // Initialize structured logging
    logger = log.New(os.Stdout, "[STARTTECH] ", log.LstdFlags|log.Lshortfile)
    
    // Initialize Redis
    redisClient = redis.NewClient(&redis.Options{
        Addr:     getEnv("REDIS_ENDPOINT", "localhost:6379"),
        Password: "",
        DB:       0,
    })
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := redisClient.Ping(ctx).Err(); err != nil {
        logger.Printf("Warning: Redis connection failed: %v", err)
    }
    
    // Initialize MongoDB
    mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
    mongoClient, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    
    // Initialize CloudWatch
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(getEnv("AWS_REGION", "us-east-1")),
    }))
    cloudwatchClient := cloudwatchlogs.New(sess)
    
    // Setup router
    r := mux.NewRouter()
    
    r.HandleFunc("/health", healthHandler).Methods("GET")
    r.HandleFunc("/api/items", getItemsHandler).Methods("GET")
    r.HandleFunc("/api/items", createItemHandler).Methods("POST")
    
    // Middleware
    r.Use(loggingMiddleware)
    r.Use(corsMiddleware)
    
    // Graceful shutdown
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    go func() {
        logger.Printf("Server starting on :8080")
        if err := srv.ListenAndServe(); err != nil {
            logger.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    
    ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    _ = cloudwatchClient
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    services := make(map[string]string)
    
    // Check Redis
    if err := redisClient.Ping(r.Context()).Err(); err != nil {
        services["redis"] = "unhealthy"
    } else {
        services["redis"] = "healthy"
    }
    
    // Check MongoDB
    if err := mongoClient.Ping(r.Context(), nil); err != nil {
        services["mongodb"] = "unhealthy"
    } else {
        services["mongodb"] = "healthy"
    }
    
    overallStatus := "healthy"
    for _, status := range services {
        if status != "healthy" {
            overallStatus = "degraded"
            break
        }
    }
    
    resp := HealthResponse{
        Status:    overallStatus,
        Timestamp: time.Now(),
        Services:  services,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        logger.Printf("%s %s %v", r.Method, r.RequestURI, time.Since(start))
    })
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
