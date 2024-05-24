package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RESTEndpoint struct {
	urlpath string
	verb    string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[RESTEndpoint]http.HandlerFunc
	Middlewares   []func(next http.Handler) http.Handler
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	if string(serverPort[0]) != ":" {
		serverPort = ":" + serverPort
	}
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[RESTEndpoint]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(urlpath string, verb string, handler http.HandlerFunc) {
	s.Handlers[RESTEndpoint{
		urlpath: urlpath,
		verb:    verb,
	}] = handler
}
func (s *WebServer) AddMiddleware(midl func(next http.Handler) http.Handler) {
	s.Middlewares = append(s.Middlewares, midl)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)
	for _, handler := range s.Middlewares {
		s.Router.Use(handler)
	}
	for restEndpointInfo, handler := range s.Handlers {
		urlpath := restEndpointInfo.urlpath
		switch verb := restEndpointInfo.verb; verb {
		case http.MethodGet:
			s.Router.Get(urlpath, handler)
		case http.MethodPost:
			s.Router.Post(urlpath, handler)
		case http.MethodPut:
			s.Router.Put(urlpath, handler)
		case http.MethodPatch:
			s.Router.Patch(urlpath, handler)
		case http.MethodDelete:
			s.Router.Delete(urlpath, handler)
		default:
			return errors.New("invalid HTTP Verb")
		}

	}
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
