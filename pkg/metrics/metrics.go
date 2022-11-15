package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"social/pkg/log"
	"strconv"
)

const (
	socialNamespace = "social"
)

var (
	// buckets involves durations in milliseconds,
	// [1 2 4 8 16 32 64 128 256 512 1024 2048 4096 8192 16384 32768 65536 1.31072e+05]
	buckets = prometheus.ExponentialBuckets(1, 2, 18)
)

func Register(r *prometheus.Registry) {
	register(&HTTPHandler{
		Path:    "/metrics",
		Handler: promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	})
	register(&HTTPHandler{
		Path:    "/metrics_default",
		Handler: promhttp.Handler(),
	})
}

const (
	DefaultListenPort = "9091"
	ListenPortEnvKey  = "METRICS_PORT"
)

type HTTPHandler struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Handler     http.Handler
}

func register(h *HTTPHandler) {
	if h.HandlerFunc != nil {
		http.HandleFunc(h.Path, h.HandlerFunc)
		return
	}
	if h.Handler != nil {
		http.Handle(h.Path, h.Handler)
	}
}

func ServeHTTP() {
	go func() {
		bindAddr := getHTTPAddr()
		log.Info("listen", zap.String("addr", bindAddr))
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			log.Error("handle metrics failed", zap.Error(err))
		}
	}()
}

func getHTTPAddr() string {
	port := os.Getenv(ListenPortEnvKey)
	_, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Sprintf(":%s", DefaultListenPort)
	}

	return fmt.Sprintf(":%s", port)
}
