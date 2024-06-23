package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	statusCodeSequence []int
	currentIndex       int
	mu                 sync.Mutex
)

func initConfig() {
	viper.SetDefault("logs.level", "debug")
	viper.SetDefault("logs.format", "json")
	viper.SetEnvPrefix("logs")
	viper.AutomaticEnv()
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(viper.GetString("logs.level"))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(level)

	switch viper.GetString("logs.format") {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		})
	}

	return logger
}

func initTracer() {
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(log.Writer()), stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func parseStatusCodeSequence(codeSequence string) []int {
	codes := strings.Split(codeSequence, "-")
	var statusCodeSeq []int
	for _, code := range codes {
		if num, err := strconv.Atoi(code); err == nil {
			statusCodeSeq = append(statusCodeSeq, num)
		}
	}
	return statusCodeSeq
}

func handleRequests(logger *logrus.Logger, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("echo-server")
	ctx, span := tracer.Start(ctx, "handleRequests")
	defer span.End()

	query := r.URL.Query()
	headers := r.Header

	if codeSequence, ok := query["echo_code"]; ok {
		handleStatusCodeSequence(ctx, logger, w, r, codeSequence[0])
		return
	} else if codeSequence := headers.Get("x-echo-code"); codeSequence != "" {
		handleStatusCodeSequence(ctx, logger, w, r, codeSequence)
		return
	}

	if _, ok := query["echo_time"]; ok {
		currentTime := time.Now().Format(time.RFC3339)
		fmt.Fprintf(w, "Current Time: %s", currentTime)
		logger.Info("Echoed current time")
		return
	}

	if _, ok := query["echo_env"]; ok {
		for key, value := range headers {
			fmt.Fprintf(w, "%s: %s\n", key, value)
		}
		logger.Info("Echoed request headers")
		return
	}

	if _, ok := query["echo_body"]; ok {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read the request body: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(body) == 0 {
			fmt.Fprint(w, "No body content provided")
		} else {
			fmt.Fprintf(w, "%s", body)
		}
		logger.Info("Processed request to echo body")
		return
	}

	fmt.Fprint(w, "Welcome to the Echo Server")
	logger.Info("Served the default welcome response")
}

func handleStatusCodeSequence(ctx context.Context, logger *logrus.Logger, w http.ResponseWriter, r *http.Request, codeSequence string) {
	mu.Lock()
	defer mu.Unlock()

	tracer := otel.Tracer("echo-server")
	_, span := tracer.Start(ctx, "handleStatusCodeSequence")
	defer span.End()

	if len(statusCodeSequence) == 0 || r.URL.Query().Get("init") == "1" {
		statusCodeSequence = parseStatusCodeSequence(codeSequence)
		currentIndex = 0
	}

	if len(statusCodeSequence) > 0 {
		httpCode := statusCodeSequence[currentIndex]
		w.WriteHeader(httpCode)
		fmt.Fprintf(w, "HTTP/1.1 %d\n", httpCode)
		logger.Infof("Responded with HTTP status code: %d", httpCode)

		currentIndex = (currentIndex + 1) % len(statusCodeSequence)
	} else {
		http.Error(w, "Invalid status code sequence", http.StatusBadRequest)
		logger.Error("Invalid status code sequence: ", codeSequence)
	}
}

func main() {
	initConfig()
	initTracer()
	logger := initLogger()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequests(logger, w, r)
	})
	logger.Info("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
