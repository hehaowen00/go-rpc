package rpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	pathrouter "github.com/hehaowen00/path-router"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	ErrInvalidRequest = errors.New("invalid request payload")
	ErrNoResponse     = errors.New("no response")
)

type ServiceImpl interface {
	Init(service *Service)
}

type ServiceHandler[Req any, Res any] func(
	ctx context.Context,
	req *Req,
) (*Res, error)

type Service struct {
	addr   string
	prefix string
	lookup map[string]string
	once   sync.Once
	router pathrouter.IRouter
	server *http.Server
}

func NewService(name string, addr string, service ServiceImpl) *Service {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	serv := &Service{
		prefix: name,
		router: pathrouter.NewRouter(),
		addr:   addr,
		lookup: map[string]string{},
		server: nil,
	}

	serv.router.Use(loggerMiddleware)
	serv.router.Use(pathrouter.GzipMiddleware)
	service.Init(serv)

	return serv
}

func (s *Service) Name() string {
	return s.prefix
}

func (s *Service) add(path string, handler http.HandlerFunc) {
	if path == "" {
		panic("invalid path")
	}

	endpoint, _ := url.JoinPath(s.prefix, path)
	s.lookup[path] = endpoint
	s.router.Post(endpoint, http.StripPrefix(s.prefix, handler).ServeHTTP)
}

func (s *Service) Run() {
	s.router.Post(s.prefix, func(w http.ResponseWriter, r *http.Request) {
		EncodeJSON(w, &s.lookup)
	})

	go func() {
		h2s := &http2.Server{}
		server := http.Server{
			Addr:    s.addr,
			Handler: h2c.NewHandler(s.router, h2s),
		}

		s.server = &server

		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()
}

func (s *Service) Stop() {
	if s.server != nil {
		s.once.Do(func() {
			s.server.Shutdown(context.Background())
			s.server = nil
		})
	}
}

func Register[Req any, Res any](
	service *Service,
	name string,
	handler ServiceHandler[Req, Res],
) {
	h := proc(handler)
	service.add(name, h)
}

func proc[Req any, Res any](handler ServiceHandler[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := DecodeJSON[Req](r.Body)
		if err != nil {
			JsonError(w, http.StatusBadRequest, ErrInvalidRequest)
			return
		}

		resp, err := handler(r.Context(), &req)
		if err != nil {
			JsonError(w, http.StatusInternalServerError, err)
			return
		}

		if resp == nil {
			JsonError(w, http.StatusInternalServerError, ErrNoResponse)
			return
		}

		EncodeJSON(w, &response[Res]{
			Payload: resp,
		})
	}
}

func loggerMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
}
