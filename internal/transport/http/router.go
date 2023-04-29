package http

import "net/http"

func (s *Server) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
