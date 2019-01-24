package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func newRouter(html *html, job *api) http.Handler {
	router := httprouter.New()
	router.GET("/", html.Index)
	router.POST("/", job.upload())
	router.GET("/products", job.products())
	router.GET("/versions", job.versions())
	router.GET("/attributes", job.attributes())
	router.GET("/names", job.names())
	router.GET("/measurements", job.measurements())
	router.GET("/chart", job.chart())
	router.ServeFiles("/static/*filepath", http.Dir("assets/static"))
	return recovery(noDirListing(router))
}
