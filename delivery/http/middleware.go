package http

import (
	"net/http"
	"strings"

	"github.com/zitryss/perfmon/internal/log"
)

func recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic")
				log.Error(err)
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
