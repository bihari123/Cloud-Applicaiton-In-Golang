package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bihari123/cloud-application-in-golang/loghelper"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
}

type Options struct {
	Host string
	Port int
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	mux := chi.NewMux()

	/*
			The http.Server from the standard library has timeouts for reading different parts of HTTP requests, and for writing the HTTP responses. If it takes longer than this, the request is aborted.

		By default, no timeouts are set. When dealing with network requests on the internet, it is generally always a good idea to set timeouts, and this is no exception. In that way, we can avoid slow clients (or malicious bots) taking up server resources.
	*/

	return &Server{
		address: address,
		mux:     mux,

		// read about the following package from terminal with
		//  $ go doc http.Server

		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

/*
The life of a Server

The Server is the main thing we start up when the app starts, and what we shut down when the app stops. It's the core part of the web app, which talks to all clients and sends them your beautiful web pages as responses.

When the Server starts, it should set up your HTTP routes on the mux, and start listening on the given host and port for HTTP requests. If an error occurs, for example because something is already listening on that port, it should report the error back to the calling code.

When it stops, we should stop accepting new HTTP requests, but finish existing ones. That way, users of your web app will not see weird error messages when their connection is reset abruptly, but instead never notice that your app was down in the first place. The most typical reason for stopping the web app is when deploying a new version. The
in front will switch over from the old app version to the new one, and no-one should notice. Like magic!
*/

func (s *Server) Start() error {

	s.setupRoutes()

	loghelper.Log("Starting On ", s.address)

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server closed error")
		}
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil

}

func (s *Server) Stop() error {
	loghelper.Log("Stopping")

	/*
		shutdown timeout using the context.WithTimeout function, so that our graceful shutdown has a time limit. If the limit is reached, the server does a hard exit instead, so we are sure that we can still shut down our app successfully without the OS having to kill it. The 30 seconds here should be enough for that, even if we create long-running code later.
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping the server: %w", err)
	}
	return nil

}
