package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	authentication "github.com/qredo-external/go-rnov/pkg/auth"
	"github.com/qredo-external/go-rnov/pkg/http/handler"
	"github.com/qredo-external/go-rnov/pkg/http/middleware"
	"github.com/qredo-external/go-rnov/pkg/service"
	"github.com/qredo-external/go-rnov/pkg/storage"
)

func main() {
	// note this data should be loaded from a config file and a vault like service as well as host and port
	secret := "jdnfksdmfksd"
	td := time.Minute * 60

	virtualStorage := storage.NewUserAccess()
	auth := authentication.NewAuth(secret, td)

	authMid := middleware.NewAuthMiddleware(auth, virtualStorage)

	authSrv := service.NewAuthService(auth, virtualStorage)
	opSrv := service.NewOperationsService()

	ha := handler.NewAuthHandler(authSrv)
	ho := handler.NewOperationHandler(opSrv)

	r := mux.NewRouter()
	r.HandleFunc("/auth", ha.CreateAuthHandler).Methods("POST")
	r.HandleFunc("/sum", middleware.Authentication(*authMid, ho.SumHandler)).Methods("POST")

	fmt.Println("starting server")
	// Fire up the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
