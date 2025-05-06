package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/errs"
)

type Router struct {
	mux *http.ServeMux
}

func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(errs.Error)
			if !ok {
				log.Printf("unexpected panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println(err.Err)
			w.WriteHeader(err.Code)
			w.Write([]byte(err.Message))
		}
	}()

	router.mux.ServeHTTP(w, r)
}

func (router *Router) setRoute(httpMethod, path string, handler http.Handler) {
	pattern := fmt.Sprintf("%s %s", httpMethod, path)
	router.mux.Handle(pattern, handler)
}

func (router *Router) POST(path string, handler http.Handler) {
	router.setRoute(http.MethodPost, path, handler)
}

func (router *Router) GET(path string, handler http.Handler) {
	router.setRoute(http.MethodGet, path, handler)
}

func (router *Router) PUT(path string, handler http.Handler) {
	router.setRoute(http.MethodPut, path, handler)
}

func (router *Router) PATCH(path string, handler http.Handler) {
	router.setRoute(http.MethodPatch, path, handler)
}

func (router *Router) DELETE(path string, handler http.Handler) {
	router.setRoute(http.MethodDelete, path, handler)
}
