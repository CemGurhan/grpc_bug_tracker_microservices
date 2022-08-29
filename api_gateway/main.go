package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	mw "github.com/cemgurhan/api_gateway/authMiddleware"
	gen "github.com/cemgurhan/api_gateway/buf/gen/proto_files"
	impl "github.com/cemgurhan/api_gateway/bugs_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

var grpcAddress = ":9090"

func grpcServer() {

	grpcLis, listenErr := net.Listen("tcp", grpcAddress)

	if listenErr != nil {
		log.Fatalln("Failed to listen: ", listenErr)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterBugServiceServer(grpcServer, &impl.BugServiceServer{})

	log.Printf("serving gRPC server on port %v \n", grpcAddress)
	go func() {
		log.Fatalln(grpcServer.Serve(grpcLis))
	}()

}

func apiGatewayProxyServer() {

	conn, clientConnErr := grpc.DialContext(
		context.Background(),
		grpcAddress,
		grpc.WithBlock(),    // running locally so no need to worry about security
		grpc.WithInsecure(), // running locally so no need to worry about security
	)

	if clientConnErr != nil {
		log.Fatalln("Failed to dial grpc server: ", clientConnErr)
	}

	gwmux := runtime.NewServeMux()

	registerBugServiceErr := gen.RegisterBugServiceHandler(
		context.Background(),
		gwmux,
		conn,
	)

	if registerBugServiceErr != nil {
		log.Fatalln("Failed to register gateway:", registerBugServiceErr)
	}

	gatewayPort := ":8081"

	gwServer := &http.Server{
		Addr:    gatewayPort,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC gateway on port %v \n", gatewayPort)

	log.Fatalln(gwServer.ListenAndServe())

}

func HomePage(w http.ResponseWriter, r *http.Request) {

	grpcServer()

	apiGatewayProxyServer()

}

func main() {

	http.Handle("/", mw.IsAuthorized(HomePage))
	log.Fatal(http.ListenAndServe(":9001", nil))

}
