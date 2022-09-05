package main

import (
	"context"
	"fmt"
	"net/http"

	usercontext "github.com/cemgurhan/auth-microservice/context"
	mw "github.com/cemgurhan/auth-microservice/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/pflag"
)

var addr string

func init() {

	pflag.StringVarP(&addr, "address", "a", ":9001", "the address for the API to listen to")
	pflag.Parse()

}

type Server struct {
	Router *chi.Mux
}

func RequestID(ctx context.Context) string {
	requestID := ctx.Value(usercontext.UserContextKey)

	if requestID == nil {
		return "none"
	}

	return requestID.(string)

}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {

	s.Router.Use(mw.IsAuthorized)
	s.Router.Get("/auth", getUserFromContext)

}

func getUserFromContext(w http.ResponseWriter, r *http.Request) {

	contextuser := usercontext.UserFromContext(r.Context())

	w.Write([]byte(fmt.Sprintf("user: %v", contextuser)))

}

func main() {

	s := CreateNewServer()

	s.MountHandlers()

	fmt.Println("Listening on port:", addr)

	http.ListenAndServe(addr, s.Router)

}
