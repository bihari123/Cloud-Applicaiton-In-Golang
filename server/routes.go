package server 

import "github.com/bihari123/cloud-application-in-golang/handlers"

func (s *Server)setupRoutes(){

  handlers.Health(s.mux)
}
