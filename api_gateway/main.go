package main

import (
	"context"
	"log"
	"net"
	"net/http"

	gen "github.com/cemgurhan/api_gateway/buf/gen/proto_files"
	impl "github.com/cemgurhan/api_gateway/bugs_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// var (
// 	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
// )

// func runHttpServer() error {

// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	fmt.Printf("serving http proxy on port %v \n", 8081)
// 	fmt.Printf("serving grpc server on port %v \n", 9090)
// 	//grpc server endpoint
// 	grpcMux := runtime.NewServeMux()
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
// 	err := gen.RegisterBugServiceHandlerFromEndpoint(ctx, grpcMux, "localhost:9090", opts)
// 	if err != nil {
// 		return err
// 	}

// 	// start http server, and proxy calls to grpc server endpoint
// 	return http.ListenAndServe(":8081", grpcMux)

// }

// func main() {

// 	// flag.Parse()

// 	if err := runHttpServer(); err != nil {
// 		log.Fatal(err)
// 	}

// 	grpc.Dial(":8080")

// }

func main() {

	grpcAddress := ":9090"
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

	// ------- //

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
	log.Fatal(gwServer.ListenAndServe())

}
