package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(path string, port int) error {
	http.Handle(path, promhttp.Handler())

	addStr := fmt.Sprintf(":%d", port)

	return http.ListenAndServe(addStr, nil)
}
