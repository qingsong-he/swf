package swf

import "net/http"

func NewSwf(addr string) *swf {
	mux := http.NewServeMux()
	return &swf{
		mux: mux,
		srv: &http.Server{Addr: addr, Handler: mux},
	}
}

type swf struct {
	mux *http.ServeMux
	srv *http.Server
}

func (s *swf) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), mids ...HttpMiddleware) {
	s.mux.HandleFunc(pattern, NewHTTPMiddlewareChain(mids...)(handler))
}

func (s *swf) Run() error {
	return s.srv.ListenAndServe()
}

func (s *swf) RunWithTLS(certFile, keyFile string) error {
	return s.srv.ListenAndServeTLS(certFile, keyFile)
}

func (s *swf) Stop() error {
	return s.srv.Close()
}

// type 'HttpMiddleware' full reference from https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
type HttpMiddleware func(http.HandlerFunc) http.HandlerFunc

func NewHTTPMiddlewareChain(mw ...HttpMiddleware) HttpMiddleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}
